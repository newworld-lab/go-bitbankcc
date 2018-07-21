package entity

type Depth struct {
	Asks      [][]float64 `json:"asks,string"`
	Bids      [][]float64 `json:"bids,string"`
	Timestamp int         `json:"timestamp"`
}
