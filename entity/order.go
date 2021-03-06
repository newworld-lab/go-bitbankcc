package entity

import (
	"time"
)

type TypeOrderStatus string

const (
	OrderStatusUnfilled                TypeOrderStatus = "UNFILLED"
	OrderStatusPartiallyFilled         TypeOrderStatus = "PARTIALLY_FILLED"
	OrderStatusFullyFilled             TypeOrderStatus = "FULLY_FILLED"
	OrderStatusCanceledUnfilled        TypeOrderStatus = "CANCELED_UNFILLED"
	OrderStatusCanceledPartiallyFilled TypeOrderStatus = "CANCELED_PARTIALLY_FILLED"
)

type TypeOrderSide string

const (
	OrderSideBuy  TypeOrderSide = "buy"
	OrderSideSell TypeOrderSide = "sell"
)

type TypeOrderType string

const (
	OrderTypeLimit  TypeOrderType = "limit"
	OrderTypeMarket TypeOrderType = "market"
)

type Orders []Order

type Order struct {
	OrderID         int             `json:"order_id"`
	Pair            TypePair        `json:"pair"`
	Side            TypeOrderSide   `json:"side"`
	Type            TypeOrderType   `json:"type"`
	StartAmount     float64         `json:"start_amount,string"`
	RemainingAmount float64         `json:"remaining_amount,string"`
	ExecutedAmount  float64         `json:"executed_amount,string"`
	Price           float64         `json:"price,string"`
	AveragePrice    float64         `json:"average_price,string"`
	OrderedAt       time.Time       `json:"ordered_at"`
	ExecutedAt      *time.Time      `json:"executed_at,omitempty"`
	Status          TypeOrderStatus `json:"status"`
}

type GetOrderParams struct {
	Pair    TypePair `json:"pair"`
	OrderID int      `json:"order_id"`
}

type GetActiveOrdersParams struct {
	Pair   TypePair   `json:"pair"`
	Count  float64    `json:"count"`
	FromID string     `json:"from_id"`
	EndID  string     `json:"to_id"`
	Since  *time.Time `json:"since"`
	End    *time.Time `json:"end"`
}

type PostOrderParams struct {
	Pair   TypePair      `json:"pair"`
	Amount float64       `json:"amount,string"`
	Price  float64       `json:"price"`
	Side   TypeOrderSide `json:"side"`
	Type   TypeOrderType `json:"type"`
}

type PostCancelOrderParams struct {
	Pair    TypePair `json:"pair"`
	OrderID int      `json:"order_id"`
}

type PostCancelOrdersParams struct {
	Pair     TypePair `json:"pair"`
	OrderIDs []int    `json:"order_ids"`
}

type PostOrdersInfoParams struct {
	Pair     TypePair `json:"pair"`
	OrderIDs []int    `json:"order_ids"`
}
