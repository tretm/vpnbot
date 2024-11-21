package models

import "time"

type IpStory struct {
	Id        int       `json:"id"`
	Ip        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Country   string    `json:"country"`
	City      string    `json:"city"`
	Provider  string    `json:"provider"`
	Company   string    `json:"company"`
	Link      string    `json:"link"`
	Date      time.Time `json:"date"`
}
