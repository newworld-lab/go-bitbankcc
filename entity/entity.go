package entity

import "time"

type Ticker struct {
	Sell float64 `json:"sell,string"`
	Buy  float64 `json:"buy,string"`
	High float64 `json:"high,string"`
	Low  float64 `json:"low,string"`
	Last float64 `json:"last,string"`
	Vol  float64 `json:"vol,string"`
}

type Depth struct {
	Asks [][]float64 `json:"asks,string"`
	Bids [][]float64 `json:"bids,string"`
}

type Transaction struct {
	TransactionId int       `json:"transaction_id"`
	Side          string    `json:"side"`
	Price         float64   `json:"price,string"`
	Amount        float64   `json:"amount,string"`
	ExecutedAt    time.Time `json:"executed_at"`
}
