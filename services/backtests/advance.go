package backtests

import (
	"context"
	"fmt"
	"log"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/backtest"
	"github.com/lerenn/cryptellation/pkg/event"
)

func (b Backtests) Advance(ctx context.Context, backtestId uint) error {
	return b.db.LockedBacktest(ctx, backtestId, func(bt *backtest.Backtest) (err error) {
		// Advance backtest
		finished := bt.Advance()
		log.Printf("Advancing backtest %d: %s", backtestId, bt.CurrentTime())

		// Get actual events
		evts := make([]event.Event, 0, 1)
		if !finished {
			evts, err = b.readActualEvents(ctx, *bt)
			if err != nil {
				return fmt.Errorf("cannot read actual events: %w", err)
			}
			if len(evts) == 0 {
				log.Println("WARNING: no event detected for", bt.CurrentCsTick.Time)
				bt.SetCurrentTime(bt.EndTime)
				finished = true
			} else if !evts[0].Time.Equal(bt.CurrentCsTick.Time) {
				log.Println("WARNING: no event between", bt.CurrentCsTick.Time, "and", evts[0].Time)
				bt.SetCurrentTime(evts[0].Time)
			}
		}

		// Add backtest status event
		evts = append(evts, event.NewStatusEvent(bt.CurrentCsTick.Time, event.Status{
			Finished: finished,
		}))
		b.broadcastEvents(backtestId, evts)

		return nil
	})
}

func (b Backtests) readActualEvents(ctx context.Context, bt backtest.Backtest) ([]event.Event, error) {
	evts := make([]event.Event, 0, len(bt.TickSubscriptions))
	for _, sub := range bt.TickSubscriptions {
		list, err := b.candlesticks.Read(ctx, client.ReadCandlesticksPayload{
			ExchangeName: sub.ExchangeName,
			PairSymbol:   sub.PairSymbol,
			Period:       bt.PeriodBetweenEvents,
			Start:        &bt.CurrentCsTick.Time,
			End:          &bt.EndTime,
			Limit:        1,
		})
		if err != nil {
			return nil, fmt.Errorf("could not get candlesticks from service: %w", err)
		}

		tcs, exists := list.First()
		if !exists {
			continue
		}

		evt, err := event.TickEventFromCandlestick(sub.ExchangeName, sub.PairSymbol, bt.CurrentCsTick.PriceType, tcs.Time, tcs.Candlestick)
		if err != nil {
			return nil, fmt.Errorf("turning candlestick into event: %w", err)
		}
		evts = append(evts, evt)
	}

	_, evts = event.OnlyKeepEarliestSameTimeEvents(evts, bt.EndTime)
	log.Printf("%d events for ticks on backtest %d", len(evts), bt.ID)
	return evts, nil
}

func (b Backtests) broadcastEvents(backtestId uint, evts []event.Event) {
	log.Printf("Broadcasting %d events on backtest %d", len(evts), backtestId)

	var count uint
	for _, evt := range evts {
		log.Printf("Broadcasting event %+v for backtest %d", evt, backtestId)
		if err := b.events.Publish(backtestId, evt); err != nil {
			log.Println("WARNING: error when publishing event", evt)
			continue
		}

		count++
	}

	if count == 0 {
		log.Println("WARNING: no available events")
	}
}
