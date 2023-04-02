// Generate code for mock
//go:generate go run github.com/golang/mock/mockgen -source=port.go -destination=mock.gen.go -package exchanges

package exchanges

import (
	"context"

	"github.com/digital-feather/cryptellation/pkg/exchange"
)

type Port interface {
	Infos(ctx context.Context, name string) (exchange.Exchange, error)
}
