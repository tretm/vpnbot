package vpn

import (
	"vpngigabot/internal/config"
	"vpngigabot/internal/models"
	"vpngigabot/internal/vpn/outline"
	"vpngigabot/internal/vpn/vtoray"
)

const (
	Outline string = "outline"
	VtoRay  string = "vtoray"
)

type VPNManager interface {
	CreateKey(*models.User) (*models.Link, error)
	SetTraficLimit(key string, limitGb int) error
	DeleteKey(string) error
	UpdateKey(*models.Link) error
	KeyBuilder(*models.Link) string
}

func New(typeVPN string, cfg *config.Config) VPNManager {
	switch typeVPN {
	case Outline:
		return outline.New(cfg)
	case VtoRay:
		return vtoray.New(cfg)
	default:
		return outline.New(cfg)
	}
}
