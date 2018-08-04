package entity

type Accounts []Account

type Account struct {
	Uuid    string `json:"uuid"`
	Label   string `json:"label"`
	Address string `json:"address"`
}
