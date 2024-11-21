package messages

import (
	"fmt"
	"log"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	BalanceAdd = "addmoney"
)

func (mm *Messages) succesPay(user *models.User) tgbotapi.MessageConfig {
	const op = "messages.succesPay"
	// bot.DeleteMessage(tgbotapi.NewDeleteMessage(p.user.Userid, p.user.Messageid))
	var keys [][]tgbotapi.InlineKeyboardButton
	keys = append(keys,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmOkText, BtmCancelValue)),
	)
	msg := tgbotapi.NewMessage(user.UserId, fmt.Sprintf(MsgPayAdded, user.Balance, user.BalanceAllTime)) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)

	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	msg.ReplyMarkup = &numericKeyboard1
	phs := mm.sm.NewPaymentHistoryStorage()
	tx, _, err := phs.Insert(nil, &models.PaymentHistory{UserId: user.UserId, Amount: user.Balance, TransactionType: BalanceAdd})
	if err != nil {
		tx.Rollback()
		log.Println(op, err)
	}
	err = tx.Commit()
	if err != nil {
		log.Println(op, err)
	}
	user.Balance = 0
	return msg
}
