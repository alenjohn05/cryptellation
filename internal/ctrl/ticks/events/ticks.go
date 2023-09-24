// Ticks
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.24.3 -g application -p events -i ./../../../../api/asyncapi/ticks.yaml -o ./app.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.24.3 -g user        -p events -i ./../../../../api/asyncapi/ticks.yaml -o ./user.gen.go
//go:generate go run github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@v0.24.3 -g types       -p events -i ./../../../../api/asyncapi/ticks.yaml -o ./types.gen.go

package events

import (
	"time"

	client "github.com/lerenn/cryptellation/clients/go"
	"github.com/lerenn/cryptellation/pkg/models/tick"
)

func (msg *RegisteringRequestMessage) Set(payload client.TicksFilterPayload) {
	msg.Payload.Exchange = ExchangeNameSchema(payload.ExchangeName)
	msg.Payload.Pair = PairSymbolSchema(payload.PairSymbol)
}

func (msg *TickMessage) ToModel() tick.Tick {
	return tick.Tick{
		Time:       time.Time(msg.Payload.Time),
		PairSymbol: string(msg.Payload.PairSymbol),
		Price:      msg.Payload.Price,
		Exchange:   string(msg.Payload.Exchange),
	}
}
