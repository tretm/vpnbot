package messages

import (
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) instruction(user *models.User) tgbotapi.MessageConfig {

	var keys [][]tgbotapi.InlineKeyboardButton
	text := ""
	switch user.GetCommandData() {
	case BtmAndroidValue:
		text = MsgAndroidText
	case BtmIOSValue:
		text = MsgIOSText
	case BtmWindowsValue:
		text = MsgWindowsText
	case BtmLinuxText:
		text = MsgLinuxText
	default:
		text = MsgAndroidText
	}

	keys = append(keys,

		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmCancelText, BtmHelpValue)),
	)
	msg := tgbotapi.NewMessage(user.UserId, text) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	msg.ReplyMarkup = &numericKeyboard1
	return msg
}
