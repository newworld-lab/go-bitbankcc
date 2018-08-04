package entity

import "time"

type Trades []Trade

type TypeTradeOrder string

const (
	Asc  TypeTradeOrder = "asc"
	Desc TypeTradeOrder = "desc"
)

type Trade struct {
	TradeId        int    `json:"trade_id"`
	Pair           string `json:"pair"`
	OrderId        int    `json:"order_id"`
	Side           string `json:"side"`
	Type           string `json:"type"`
	Amount         string `json:"amount"`
	Price          string `json:"price"`
	MakerTaker     string `json:"maker_taker"`
	FeeAmountBase  string `json:"fee_amount_base"`
	FeeAmountQuote string `json:"fee_amount_quote"`
	ExecuteAt      int    `json:"executed_at"`
}

type TradeParams struct {
	Pair    TypePair       `json:"pair"`
	Count   float64        `json:"count"`
	OrderId float64        `json:"order_id"`
	Since   *time.Time     `json:"since"`
	End     *time.Time     `json:"end"`
	Order   TypeTradeOrder `json:"order"`
}
