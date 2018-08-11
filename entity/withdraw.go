package entity

type Accounts []Account

type Account struct {
	UUID    string `json:"uuid"`
	Label   string `json:"label"`
	Address string `json:"address"`
}
