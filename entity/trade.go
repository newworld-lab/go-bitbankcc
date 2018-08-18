package entity

import "time"

type Trades []Trade

type TypeTradeOrder string

const (
	Asc  TypeTradeOrder = "asc"
	Desc TypeTradeOrder = "desc"
)

type Trade struct {
	TradeId        int       `json:"trade_id"`
	Pair           string    `json:"pair"`
	OrderId        int       `json:"order_id"`
	Side           float64   `json:"side,string"`
	Type           string    `json:"type"`
	Amount         float64   `json:"amount,string"`
	Price          float64   `json:"price,string"`
	MakerTaker     string    `json:"maker_taker"`
	FeeAmountBase  float64   `json:"fee_amount_base,string"`
	FeeAmountQuote float64   `json:"fee_amount_quote,string"`
	ExecutedAt     time.Time `json:"executed_at"`
}

type TradeParams struct {
	Pair    TypePair       `json:"pair"`
	Count   float64        `json:"count"`
	OrderId float64        `json:"order_id"`
	Since   *time.Time     `json:"since"`
	End     *time.Time     `json:"end"`
	Order   TypeTradeOrder `json:"order"`
}
