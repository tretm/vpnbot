package models

import (
	"time"
)

const (
	StateAllowed = "ok"
	StateBanned  = "ban"
)

type Link struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
	Link        string    `json:"link"`
	VpnLink     string    `json:"vpn_link"`
	VpnLinkId   string    `json:"vpn_link_id"`
	VpnPassword string    `json:"vpn_lint_password"`
	State       string    `json:"state"`
	TimeEnd     time.Time `json:"time_end"`
	TimeCreate  time.Time `json:"time_create"`
	TimeUpdate  time.Time `json:"time_update"`
}
