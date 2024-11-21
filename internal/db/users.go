package db

import (
	"database/sql"
	"vpngigabot/internal/models"
)

type Users interface {
	FindOne(userId int64) (*models.User, error)
	Find(filter *UserFilter, orderFilter *OrderByUsers, offset, limit int) ([]*models.User, error)
	Insert(*sql.Tx, *models.User) (*sql.Tx, int64, error)
	Update(*sql.Tx, *models.User, *UserFilter) (*sql.Tx, int64, error)
	UpdateUser(*sql.Tx, *models.User, int64) (*sql.Tx, error)
	Delete(*sql.Tx, *UserFilter) (*sql.Tx, error)
}

type UserFilter struct {
	Id               string
	UserId           string
	UserName         string
	Password         string
	Count            string
	MessageType      string
	Command          string
	Lang             string
	HistoryUserName  string
	WhateDescription string
	LinkDescription  string
	City             string
	Phone            string
	IsDeanon         string
	Status           string
	Role             string
}

type OrderByUsers struct {
	Desc       bool
	Id         bool
	UserId     bool
	Count      bool
	TimeCreate bool
	TimeUpdate bool
	TimeDeanon bool
}
