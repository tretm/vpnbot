package youmoney

import (
	"fmt"
	"vpngigabot/internal/config"
)

const baseURL = "https://yoomoney.ru/api/"

type YouMoneyClient struct {
	token     string
	account   string
	domenName string
}

func New(cfg *config.Config) (*YouMoneyClient, error) {
	if cfg.YouMoneyConfig.Token == "" {
		auth := authorize{
			ClientID:    cfg.ClientId,
			RedirectURI: cfg.RedirectUrl,
			Scope: []string{"account-info",
				"operation-history",
				"operation-details",
				"incoming-transfers",
				"payment-p2p",
				"payment-shop"},
		}

		if t, err := auth.authorize(); err != nil {
			fmt.Println("Error:", err)
			return nil, err
		} else {
			cfg.YouMoneyConfig.Token = t
		}
	}

	ymc := YouMoneyClient{token: cfg.YouMoneyConfig.Token, domenName: cfg.RedirectUrl}

	return &ymc, ymc.getAccount()
}
