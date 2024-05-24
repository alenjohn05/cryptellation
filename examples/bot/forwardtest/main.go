package main

import (
	"context"

	cryptellation "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/examples/bot"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry/console"
	"github.com/lerenn/cryptellation/pkg/adapters/telemetry/otel"
	"github.com/lerenn/cryptellation/pkg/config"
	"github.com/lerenn/cryptellation/pkg/models/account"
	"github.com/lerenn/cryptellation/svc/forwardtests/pkg/forwardtest"
)

func main() {
	// Init opentelemetry and set it globally
	console.Fallback(otel.NewTelemeter(context.Background(), "cryptellation-backtests"))
	defer telemetry.Close(context.TODO())

	// Create a cryptellation client
	client, err := cryptellation.NewServices(config.LoadNATS())
	if err != nil {
		panic(err)
	}

	// Create forwardtest
	b, err := cryptellation.NewForwardTest(client, forwardtest.NewPayload{
		Accounts: map[string]account.Account{
			"binance": {
				Balances: map[string]float64{
					"USDT": 1000,
				},
			},
		},
	}, &bot.Bot{})
	if err != nil {
		panic(err)
	}

	// Run backtest
	if err := b.Run(); err != nil {
		panic(err)
	}
}
