package candlesticks

import (
	"context"
	"log"
	"time"

	"github.com/lerenn/cryptellation/internal/core/candlesticks/ports/exchanges"
	"github.com/lerenn/cryptellation/pkg/models/candlestick"
)

const (
	// This is the minimal quantity of candlesticks that will be retrieved in case of miss
	// It will avoid too many request on exchanges if few candlesticks are requested regularly.
	MinimalRetrievedMissingCandlesticks = 100
)

func (app candlesticks) GetCached(ctx context.Context, payload GetCachedPayload) (*candlestick.List, error) {
	log.Printf("Get candlesticks for %+v", payload)

	start, end := payload.Period.RoundInterval(payload.Start, payload.End)
	cl := candlestick.NewEmptyList(payload.ExchangeName, payload.PairSymbol, payload.Period)

	// Read candlesticks from database
	if err := app.db.ReadCandlesticks(ctx, cl, start, end, payload.Limit); err != nil {
		return nil, err
	}
	log.Printf("Read %d candlesticks from %s to %s (limit: %d)", cl.Len(), start, end, payload.Limit)

	if !cl.AreMissing(start, end, payload.Limit) {
		log.Printf("No candlestick missing, returning the list with %d candlesticks.", cl.Len())
		return cl, nil
	}

	downloadStart, downloadEnd := getDownloadStartEndTimes(cl, start, end)
	if err := app.download(ctx, cl, downloadStart, downloadEnd, payload.Limit); err != nil {
		return nil, err
	}

	if err := app.upsert(ctx, cl); err != nil {
		return nil, err
	}

	return cl.Extract(start, end, payload.Limit), nil
}

// getDownloadStartEndTimes gives start and end time for download
// Time order: start < end
func getDownloadStartEndTimes(cl *candlestick.List, end, start time.Time) (time.Time, time.Time) {
	t, _, exists := cl.TimeSerie.Last()
	if exists && !cl.HasUncomplete() {
		end = t.Add(cl.Period.Duration())
	}

	qty := int(cl.Period.CountBetweenTimes(end, start)) + 1
	if qty < MinimalRetrievedMissingCandlesticks {
		d := cl.Period.Duration() * time.Duration(MinimalRetrievedMissingCandlesticks-qty)
		start = start.Add(d)
	}

	return end, start
}

func (app candlesticks) download(ctx context.Context, cl *candlestick.List, start, end time.Time, limit uint) error {
	payload := exchanges.GetCandlesticksPayload{
		Exchange:   cl.ExchangeName,
		PairSymbol: cl.PairSymbol,
		Period:     cl.Period,
		Start:      start,
		End:        end,
	}

	for {
		ncl, err := app.exchanges.GetCandlesticks(ctx, payload)
		if err != nil {
			return err
		}

		if err := cl.Merge(ncl, nil); err != nil {
			return err
		}

		cl.ReplaceUncomplete(ncl)

		t, _, exists := ncl.Last()
		if !exists || !t.Before(end) || (limit != 0 && cl.Len() >= int(limit)) {
			break
		}

		payload.Start = t.Add(cl.Period.Duration())
	}

	return nil
}

func (app candlesticks) upsert(ctx context.Context, cl *candlestick.List) error {
	tStart, _, startExists := cl.First()
	tEnd, _, endExists := cl.Last()
	if !startExists || !endExists {
		return nil
	}

	rcl := candlestick.NewEmptyListFrom(cl)
	if err := app.db.ReadCandlesticks(ctx, rcl, tStart, tEnd, 0); err != nil {
		return err
	}

	csToInsert := candlestick.NewEmptyListFrom(cl)
	csToUpdate := candlestick.NewEmptyListFrom(cl)
	if err := cl.Loop(func(ts time.Time, cs candlestick.Candlestick) (bool, error) {
		_, exists := rcl.Get(ts)
		if !exists {
			return false, csToInsert.Set(ts, cs)
		} else {
			return false, csToUpdate.Set(ts, cs)
		}
	}); err != nil {
		return err
	}

	if csToInsert.Len() > 0 {
		return app.db.CreateCandlesticks(ctx, csToInsert)
	}

	if csToUpdate.Len() > 0 {
		return app.db.UpdateCandlesticks(ctx, csToUpdate)
	}

	return nil
}
