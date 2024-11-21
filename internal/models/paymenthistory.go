package models

import "time"

type PaymentHistory struct {
	Id              int64     `json:"id"`
	UserId          int64     `json:"user_id"`
	Amount          int       `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	Comment         string    `json:"comment"`
	TimeCreate      time.Time `json:"time_create"`
}
