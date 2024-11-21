package messages

import (
	"log"
	"strings"
	"vpngigabot/internal/config"
	"vpngigabot/internal/models"
	"vpngigabot/internal/paysystem"
	tgbot "vpngigabot/internal/tgbot"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Messages struct {
	cfg *config.Config
	sm  tgbot.SotrageManager
	vpn VPNManager
	ps  paysystem.PaySystem
}

type VPNManager interface {
	CreateKey(*models.User) (*models.Link, error)
	SetTraficLimit(key string, limitGb int) error
	DeleteKey(string) error
	UpdateKey(*models.Link) error
	KeyBuilder(*models.Link) string
}

func New(cfg *config.Config, vpnMng VPNManager, sm tgbot.SotrageManager, ps paysystem.PaySystem) *Messages {
	return &Messages{cfg: cfg, vpn: vpnMng, sm: sm, ps: ps}
}

func (mm *Messages) MessageManager(user *models.User) tgbotapi.Chattable {
	const op = "messages.MessageManager"

	if strings.Contains(user.Command, "_") {
		data := strings.Split(user.Command, "_")
		if len(data) >= 2 {
			user.Command = data[0]
			user.SetCommandData(data[1])
		}
	}
	log.Printf("%s ******** UserId: %d, Command: %s ", op, user.UserId, user.Command)
	switch user.Command {
	case "start":
		return mm.start(user)
	case BtmVpnTrialValue:
		return mm.trial(user)
	case BtmGenerateKeyValue:
		return mm.newTrialOutlineKey(user)
	case BtmBalancAddValue:
		return mm.payMenu(user)
	case BtmPayMenuValue:
		return mm.generatePayLink(user)
	case tgbot.SuccesPay:
		return mm.succesPay(user)
	case BtmVpnKeyValue:
		return mm.keyMenu(user)
	case BtmVpnKeyDetaleValue:
		return mm.keyDetalisation(user)
	case BtmAddKeyValue:
		return mm.priceMeny(user)
	case BtmHistoryPayValue:
		return mm.payHistory(user)
	case BtmExtndValue:
		return mm.priceExtendMeny(user)
	case BtmPiceExtendValue, BtmPriceValue:
		return mm.confirmBuy(user)
	case BtmConfirmYesValue:
		return mm.confirmExtBuy(user)
	case BtmPayCheckValue:
		return mm.chackPay(user)
	case "deletekey":
		return mm.deleteKey(user)
	case BtmCancelValue:
		return mm.start(user)
	case BtmHelpValue, BtmHowToUseValue:
		return mm.help(user)
	case BtmOSValue:
		return mm.instruction(user)
	case BtmNotificationValue:
		return mm.notification(user)
	default:
		return mm.start(user)

	}

}
