package messages

import (
	"fmt"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) priceExtendMeny(user *models.User) tgbotapi.MessageConfig {
	const op = "messages.priceExtendMeny"

	var keys [][]tgbotapi.InlineKeyboardButton
	if len(mm.cfg.PeriodsText) != len(mm.cfg.PeriodsVal) {
		return mm.start(user)
	}
	for i, p := range mm.cfg.PeriodsText {

		discount := mm.cfg.Discount * i

		price := mm.cfg.PeriodsVal[i] * (mm.cfg.PriceOneMonth - (mm.cfg.PriceOneMonth * discount / 100))
		btmText := ""
		if discount > 0 {
			btmText = fmt.Sprintf("(-%d%%)", discount)
		}

		keys = append(keys,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf(BtmPriceText, p+btmText, price),
					BtmPiceExtendValue+fmt.Sprintf("_%d-%d-%s", mm.cfg.PeriodsVal[i], price, user.GetCommandData())),

				// tgbotapi.NewInlineKeyboardRow(
				// 	tgbotapi.NewInlineKeyboardButtonData(BtmCancelText, BtmCancelValue)),
			))
	}
	keys = append(keys,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BtmCancelText, BtmVpnKeyValue)),
	)
	msg := tgbotapi.NewMessage(user.UserId, fmt.Sprintf(MsgPriceExtendMenu, user.BalanceAllTime)) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	msg.ReplyMarkup = &numericKeyboard1
	return msg
}
