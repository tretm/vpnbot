package db

import (
	"database/sql"
	"vpngigabot/internal/models"
)

type IpStories interface {
	Find(filter *IpStoryFilter, offset, limit int) ([]*models.IpStory, error)
	Insert(*sql.Tx, *models.IpStory) (*sql.Tx, int64, error)
	Update(*sql.Tx, *models.IpStory, *IpStoryFilter) (*sql.Tx, int64, error)
	Delete(*sql.Tx, *IpStoryFilter) (*sql.Tx, error)
}

type IpStoryFilter struct {
	Id        string
	Ip        string
	UserAgent string
	Country   string
	City      string
	Provider  string
	Company   string
	Link      string
	DateStart string
	DateEnd   string
	UserId    string
}

type OrderIpStoryes struct {
	Desk bool
	Date bool
}

type IpWithUserFileter struct {
	UserId string
	Link   string
}
