package messages

import (
	"strconv"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func makeColumns(data []int, columns int) [][]tgbotapi.InlineKeyboardButton {
	rows := len(data) / columns
	if len(data)%columns != 0 {
		rows++
	}

	result := make([][]tgbotapi.InlineKeyboardButton, rows)
	for i := range result {
		if i == rows-1 && len(data)%columns != 0 {
			result[i] = make([]tgbotapi.InlineKeyboardButton, len(data)%columns)
		} else {
			result[i] = make([]tgbotapi.InlineKeyboardButton, columns)
		}
	}
	for i, val := range data {
		row := i / columns
		col := i % columns
		amount := strconv.Itoa(val)
		result[row][col] = tgbotapi.NewInlineKeyboardButtonData(amount+" ₽", BtmPayMenuValue+"_"+amount)
	}

	return result
}
func (mm *Messages) payMenu(user *models.User) tgbotapi.MessageConfig {
	// bot.DeleteMessage(tgbotapi.NewDeleteMessage(p.user.Userid, p.user.Messageid))
	if len(mm.cfg.PayAmount) < 1 {
		return mm.start(user)
	}
	keys := makeColumns(mm.cfg.PayAmount, 3)
	keys = append(keys, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmCancelText, BtmCancelValue)))
	msg := tgbotapi.NewMessage(user.UserId, BtmPayMenuText) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	msg.ReplyMarkup = &numericKeyboard1
	return msg
}
