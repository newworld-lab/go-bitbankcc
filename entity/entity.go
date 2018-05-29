package entity

import "time"

type Ticker struct {
	Sell      float64 `json:"sell,string"`
	Buy       float64 `json:"buy,string"`
	High      float64 `json:"high,string"`
	Low       float64 `json:"low,string"`
	Last      float64 `json:"last,string"`
	Vol       float64 `json:"vol,string"`
	Timestamp int     `json:"timestamp"`
}

type Depth struct {
	Asks      [][]float64 `json:"asks,string"`
	Bids      [][]float64 `json:"bids,string"`
	Timestamp int         `json:"timestamp"`
}

type Transaction struct {
	TransactionId int       `json:"transaction_id"`
	Side          string    `json:"side"`
	Price         float64   `json:"price,string"`
	Amount        float64   `json:"amount,string"`
	ExecutedAt    time.Time `json:"executed_at"`
}

type Transactions []Transaction

type OhlcvItem struct {
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume float64   `json:"volume"`
	Date   time.Time `json:"date"`
}

type Candlestick []CandlestickItem

type CandlestickItem struct {
	Type  string      `json:"type"`
	Ohlcv []OhlcvItem `json:"ohlcv"`
}

type Assets []Asset

type Asset struct {
	Asset           string      `json:"asset"`
	AmountPrecision int         `json:"amount_precision"`
	OnhandAmount    float64     `json:"onhand_amount"`
	LockedAmount    float64     `json:"locked_amount"`
	FreeAmount      float64     `json:"free_amount"`
	WithDrawalFee   interface{} `json:"withdrawal_fee"`
}
