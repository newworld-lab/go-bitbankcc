package entity

import "time"

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
