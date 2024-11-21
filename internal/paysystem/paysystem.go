package paysystem

import "vpngigabot/internal/models"

const (
	YASUCCED   = "success"
	YAREFUSED  = "refused"
	YAPROGRESS = "in_progress"
)

type PaySystem interface {
	CreatePayLink(summAmount float64) (*models.Pay, error)
	CheckPayStatus(*models.Pay) (string, error)
}
