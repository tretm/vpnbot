package messages

import (
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) start(user *models.User) tgbotapi.MessageConfig {
	// bot.DeleteMessage(tgbotapi.NewDeleteMessage(p.user.Userid, p.user.Messageid))

	var keys [][]tgbotapi.InlineKeyboardButton

	if !user.TestUsed {
		keys = append(keys, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmVpnTrialText, BtmVpnTrialValue)))
	}
	keys = append(keys, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmVpnKeyText, BtmVpnKeyValue)),
		// tgbotapi.NewInlineKeyboardButtonData(filds.BtmExitmoney, filds.ExitMoney)),
		// tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(filds.BtmRuue, filds.Ruue),
		//
		//	tgbotapi.NewInlineKeyboardButtonData(filds.BtmUeru, filds.Ueru)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmBalancAddText, BtmBalancAddValue)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmHistoryPayText, BtmHistoryPayValue)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmHelpText, BtmHelpValue)),
		// tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmReferalText, BtmReferalValue)),
	)
	msg := tgbotapi.NewMessage(user.UserId, MESSAGESTART) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	msg.ReplyMarkup = &numericKeyboard1
	return msg
}
