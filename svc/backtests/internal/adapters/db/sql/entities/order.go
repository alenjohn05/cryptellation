package entities

import (
	"time"

	"github.com/lerenn/cryptellation/svc/backtests/pkg/order"
)

type Order struct {
	ID            uint `gorm:"primaryKey"`
	BacktestID    uint
	ExecutionTime *time.Time
	Type          string
	Exchange      string
	Pair          string
	Side          string
	Quantity      float64
	Price         float64
}

func (o Order) ToModel() (order.Order, error) {
	t := order.Type(o.Type)
	if err := t.Validate(); err != nil {
		return order.Order{}, err
	}

	s := order.Side(o.Side)
	if err := s.Validate(); err != nil {
		return order.Order{}, err
	}

	return order.Order{
		ID:            uint64(o.ID),
		ExecutionTime: o.ExecutionTime,
		Type:          t,
		Exchange:      o.Exchange,
		Pair:          o.Pair,
		Side:          s,
		Quantity:      o.Quantity,
		Price:         o.Price,
	}, nil
}

func ToOrderModels(orders []Order) ([]order.Order, error) {
	var err error
	models := make([]order.Order, len(orders))
	for i, e := range orders {
		if models[i], err = e.ToModel(); err != nil {
			return nil, err
		}
	}
	return models, nil
}

func FromOrderModels(backtestID uint, models []order.Order) []Order {
	entities := make([]Order, len(models))
	for i, m := range models {
		entities[i] = FromOrderModel(backtestID, m)
	}
	return entities
}

func FromOrderModel(backtestID uint, m order.Order) Order {
	return Order{
		ID:            uint(m.ID),
		BacktestID:    backtestID,
		ExecutionTime: m.ExecutionTime,
		Type:          m.Type.String(),
		Exchange:      m.Exchange,
		Pair:          m.Pair,
		Side:          m.Side.String(),
		Quantity:      m.Quantity,
		Price:         m.Price,
	}
}
