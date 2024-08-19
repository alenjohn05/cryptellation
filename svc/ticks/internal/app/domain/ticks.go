package domain

import (
	"context"

	"cryptellation/pkg/models/event"

	"cryptellation/svc/ticks/internal/app/ports/events"
	"cryptellation/svc/ticks/internal/app/ports/exchanges"
)

type Ticks struct {
	adapters   adapters
	listenings listenings
}

func New(exchanges exchanges.Port, events events.Port) *Ticks {
	if exchanges == nil {
		panic("nil exchanges")
	}

	if events == nil {
		panic("nil events")
	}

	a := adapters{
		Exchanges: exchanges,
		Events:    events,
	}

	return &Ticks{
		adapters:   a,
		listenings: newListenings(a),
	}
}

func (t *Ticks) ListeningNotificationReceived(ctx context.Context, ts event.TickSubscription) {
	t.listenings.UpdateLastNotificationSeen(ts)
}

func (t *Ticks) Close(ctx context.Context) {
	t.listenings.Close(ctx)
}
