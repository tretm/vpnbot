package db

import (
	"database/sql"
	"vpngigabot/internal/models"
)

type Links interface {
	Find(filter *LinkFilter, order *OrderLinks, offset, limit int) ([]*models.Link, error)
	Insert(*sql.Tx, *models.Link) (*sql.Tx, int64, error)
	Update(tx *sql.Tx, item *models.Link, where *LinkFilter) (*sql.Tx, int64, error)
	Delete(*sql.Tx, *LinkFilter) (*sql.Tx, error)
}

type LinkFilter struct {
	Id              string
	UserId          string
	Link            string
	VpnLink         string
	VpnLinkId       string
	VpnLinkPassword string
	State           string
	TimeEnd         string
	TimeCreate      string
	TimeUpdate      string
}

type OrderLinks struct {
	Desk bool
	Date bool
}
