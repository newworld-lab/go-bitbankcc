package entity

import "time"

type Transaction struct {
	TransactionId int       `json:"transaction_id"`
	Side          string    `json:"side"`
	Price         float64   `json:"price,string"`
	Amount        float64   `json:"amount,string"`
	ExecutedAt    time.Time `json:"executed_at"`
}

type Transactions []Transaction
