package db

import (
	"database/sql"
	"vpngigabot/internal/models"
)

type PaymentHistory interface {
	FindOne(userId int64) (*models.PaymentHistory, error)
	Find(filter *PaymentHistoryFilter, orderFilter *OrderByPaymentHistory, offset, limit int) ([]*models.PaymentHistory, error)
	Insert(*sql.Tx, *models.PaymentHistory) (*sql.Tx, int64, error)
	Delete(*sql.Tx, *models.PaymentHistory) (*sql.Tx, error)
}

type PaymentHistoryFilter struct {
	Id              string
	UserId          string
	Amount          string
	TransactionType string
	Comment         string
	TimeCreate      string
}
type OrderByPaymentHistory struct {
	UserId          bool
	Amount          bool
	TransactionType bool
	Comment         bool
	TimeCreate      bool
}
