package db

import (
	"database/sql"
	"vpngigabot/internal/models"
)

type PayLink interface {
	Find(*PayLinkFilter) (*models.PayLink, error)
	Insert(*sql.Tx, *models.PayLink) (*sql.Tx, int64, error)
	// Update(*models.PayLink, *UserFilter) (int64, error)
	Delete(*sql.Tx, *PayLinkFilter) (*sql.Tx, error)
}

type PayLinkFilter struct {
	PayId  string
	Amount string
	UserId string
	Status string
	Date   string
}
