package models

import "time"

type PayLink struct {
	Id     int       `json:"id"`
	PayId  string    `json:"pay_id"`
	Amount int       `json:"amount"`
	UserId int64     `json:"user_id"`
	Status string    `json:"status"`
	Date   time.Time `json:"date"`
}
