package entity

type Ticker struct {
	Sell      float64 `json:"sell,string"`
	Buy       float64 `json:"buy,string"`
	High      float64 `json:"high,string"`
	Low       float64 `json:"low,string"`
	Last      float64 `json:"last,string"`
	Vol       float64 `json:"vol,string"`
	Timestamp int     `json:"timestamp"`
}
