package daemon

import (
	"context"

	client "github.com/lerenn/cryptellation/clients/go"
	natsClient "github.com/lerenn/cryptellation/clients/go/nats"
	sql "github.com/lerenn/cryptellation/internal/adapters/db/sql/indicators"
	"github.com/lerenn/cryptellation/internal/components/indicators/ports/db"
	"github.com/lerenn/cryptellation/pkg/config"
)

type adapters struct {
	db           db.Port
	candlesticks client.Candlesticks
}

func newAdapters(ctx context.Context) (adapters, error) {
	// Init database client
	db, err := sql.New(config.LoadSQL())
	if err != nil {
		return adapters{}, err
	}

	// Init candlesticks client
	candlesticks, err := natsClient.NewCandlesticks(config.LoadNATS())
	if err != nil {
		return adapters{}, err
	}

	return adapters{
		db:           db,
		candlesticks: candlesticks,
	}, nil
}

func (a adapters) Close(ctx context.Context) {
	a.candlesticks.Close(ctx)
}
