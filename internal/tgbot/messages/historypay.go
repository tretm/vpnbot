package messages

import (
	"fmt"
	"strconv"
	"vpngigabot/internal/db"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	layout = "2006-01-02"
	header = "№ |	Дата				 | Сумма"
)

func (mm *Messages) payHistory(user *models.User) tgbotapi.MessageConfig {
	const op = "messages.payHistory"
	msgText := ""
	var keys [][]tgbotapi.InlineKeyboardButton
	phs := mm.sm.NewPaymentHistoryStorage()
	strUId := strconv.FormatInt(user.UserId, 10)
	ph, err := phs.Find(&db.PaymentHistoryFilter{UserId: strUId}, &db.OrderByPaymentHistory{TimeCreate: true}, 0, mm.cfg.StorageConfig.RowsLimit)
	if err != nil {
		msgText = MsgPayHistoryEmpty
		keys = append(keys, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmOkText, BtmCancelValue)))
	} else {
		if len(ph) < 1 {
			msgText = MsgPayHistoryEmpty
			keys = append(keys, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmOkText, BtmCancelValue)))
		} else if len(ph) >= mm.cfg.StorageConfig.RowsLimit {
			msgText = fmt.Sprintf(MsgPayHistory, generateString(ph))
			keys = append(keys, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmNextPageText, BtmNextPageValue+"_"+strconv.Itoa(mm.cfg.StorageConfig.RowsLimit))),
				tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmOkText, BtmCancelValue)))
		} else {
			msgText = fmt.Sprintf(MsgPayHistory, generateString(ph))
			keys = append(keys, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmOkText, BtmCancelValue)))
		}
	}

	msg := tgbotapi.NewMessage(user.UserId, msgText) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	msg.ReplyMarkup = &numericKeyboard1
	return msg
}
func generateString(ph []*models.PaymentHistory) string {
	resStr := header + "\n"
	amount := ""
	for i, p := range ph {

		if p.TransactionType == BalanceAdd {
			amount = "+" + strconv.Itoa(p.Amount)
		} else {
			amount = "-" + strconv.Itoa(p.Amount)
		}
		resStr += strconv.Itoa(i+1) + ". " + p.TimeCreate.Format(layout) + " " + amount + "\n"

	}
	return "<pre><code>" + resStr + "</code></pre>"
}
