package entity

type Assets []Asset

type Asset struct {
	Asset           string      `json:"asset"`
	AmountPrecision int         `json:"amount_precision"`
	OnhandAmount    float64     `json:"onhand_amount"`
	LockedAmount    float64     `json:"locked_amount"`
	FreeAmount      float64     `json:"free_amount"`
	WithDrawalFee   interface{} `json:"withdrawal_fee"`
}
