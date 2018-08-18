package entity

import "time"

type Accounts []Account

type Account struct {
	UUID    string `json:"uuid"`
	Label   string `json:"label"`
	Address string `json:"address"`
}

type TypeWithdrawSide string

const (
	WithdrawSideBuy  TypeWithdrawSide = "buy"
	WithdrawSideSell TypeWithdrawSide = "sell"
)

type TypeWithdrawType string

const (
	WithdrawTypeLimit  TypeWithdrawType = "limit"
	WithdrawTypeMarket TypeWithdrawType = "market"
)

type Withdraw struct {
	UUID        string    `json:"uuid"`
	Asset       TypeAsset `json:"asset"`
	AccountUUID string    `json:"account_uuid"`
	Amount      float64   `json:"amount,string"`
	Fee         float64   `json:"fee,string"`
	Label       string    `json:"label"`
	Address     string    `json:"address"`
	Txid        string    `json:"txid"`
	Status      string    `json:"status"`
	RequestedAt time.Time `json:"requested_at"`
}

type PostWithdrawParams struct {
	Asset    TypeAsset `json:"asset"`
	UUID     string    `json:"uuid"`
	Amount   float64   `json:"amount,string"`
	OptToken string    `json:"opt_token"`
	SmsToken string    `json:"sms_token"`
}
