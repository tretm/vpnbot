package messages

import (
	"fmt"
	"log"
	"strconv"
	"vpngigabot/internal/models"
	paysistem "vpngigabot/internal/paysystem"

	//  "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) generatePayLink(user *models.User) tgbotapi.Chattable {
	const op = "generatePayLink"
	amount := 0
	pls := mm.sm.NewPayLinkStorage()
	a, err := strconv.Atoi(user.GetCommandData())
	if err != nil {
		log.Println(op, err)
		amount = 200
	}
	amount = a
	// log.Println(op, amount)
	if mm.ps != nil {

		pay, err := mm.ps.CreatePayLink(float64(amount))
		if err != nil {
			log.Println(op, err)
			return mm.payMenu(user)
		}
		tx, _, err := pls.Insert(nil, &models.PayLink{PayId: pay.PayId, Amount: amount, UserId: user.UserId, Status: paysistem.YAPROGRESS})
		if err != nil {
			tx.Rollback()
			log.Println(op, err)
			return mm.start(user)
		}
		err = tx.Commit()
		if err != nil {
			log.Println(op, err)
			return mm.start(user)
		}
		var keys [][]tgbotapi.InlineKeyboardButton
		wapp := tgbotapi.InlineKeyboardButton{Text: BtmPayUrlText, WebApp: &tgbotapi.WebAppInfo{URL: pay.PayLink}}
		// wapp := tgbotapi.InlineKeyboardButton{Text: BtmPayUrlText, WebApp: &tgbotapi.WebAppInfo{URL: "https://keysvout.tech/"}}
		keys = append(keys,
			// tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL(BtmPayUrlText, pay.PayLink)),
			tgbotapi.NewInlineKeyboardRow(wapp),
			// tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmPayCheckText, BtmPayCheckValue+"_"+pay.PayId)),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmCancelText, BtmCancelValue)),
		)
		msg := tgbotapi.NewMessage(user.UserId, fmt.Sprintf(MsgPayUrlText, amount)) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
		numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
		msg.ReplyMarkup = &numericKeyboard1
		return msg
	}
	// return tgbotapi.NewInvoice(user.UserId, MsgTitelPay, MsgDescriptionPay, "custom_payload",
	// 	mm.cfg.BotPayToken, "start_param", "RUB",
	// 	[]tgbotapi.LabeledPrice{{Label: "RUB", Amount: amount * 100}})
	return mm.start(user)

}

// create invoice
