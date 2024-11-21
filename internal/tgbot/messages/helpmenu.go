package messages

import (
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) help(user *models.User) tgbotapi.MessageConfig {

	var keys [][]tgbotapi.InlineKeyboardButton

	keys = append(keys,
		// tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmAndroidText, BtmOSValue+"_"+BtmAndroidValue)),
		// tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmIOSText, BtmOSValue+"_"+BtmIOSValue)),
		// tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmWindowsText, BtmOSValue+"_"+BtmWindowsValue)),
		// tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmLinuxText, BtmOSValue+"_"+BtmLinuxValue)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(BtmSupportText, "https://t.me/"+mm.cfg.Support)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmCancelText, BtmCancelValue)),
	)
	msg := tgbotapi.NewMessage(user.UserId, MsgHelpText) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	msg.ReplyMarkup = &numericKeyboard1
	return msg
}
