package candlesticks

import (
	"time"

	"github.com/lerenn/cryptellation/pkg/period"
)

type GetCachedPayload struct {
	ExchangeName string
	PairSymbol   string
	Period       period.Symbol
	Start        *time.Time
	End          *time.Time
	Limit        uint
}
