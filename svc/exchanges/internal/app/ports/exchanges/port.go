// Generate code for mock
//go:generate go run go.uber.org/mock/mockgen@v0.2.0 -source=port.go -destination=mock.gen.go -package exchanges

package exchanges

import (
	"context"

	"github.com/lerenn/cryptellation/exchanges/pkg/exchange"
)

type Port interface {
	Infos(ctx context.Context, name string) (exchange.Exchange, error)
}
