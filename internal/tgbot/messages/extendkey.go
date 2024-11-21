package messages

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"vpngigabot/internal/db"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) buyExtendKey(user *models.User) tgbotapi.MessageConfig {
	const op = "messages.buyExtendKey"
	var msg tgbotapi.MessageConfig
	msgText := ""
	var keys [][]tgbotapi.InlineKeyboardButton

	cd := strings.Split(user.GetCommandData(), "-")
	if len(cd) != 3 {
		log.Println(op, errors.New("Len command data less 3 "))
		return mm.start(user)
	}

	periodBuy, err := strconv.Atoi(cd[0])
	if err != nil {
		log.Println(op, err)
		return mm.start(user)
	}
	priceBuy, err := strconv.Atoi(cd[1])
	if err != nil {
		log.Println(op, err)
		return mm.start(user)
	}
	key := cd[2]
	if user.BalanceAllTime < priceBuy {
		msgText = MsgBuyErrText
		keys = append(keys,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(BtmOkText, BtmExtndValue+"_"+key)),
		)
	} else {
		ls := mm.sm.NewLinkStorage()
		tx, _, err := ls.Update(nil, &models.Link{State: models.StateAllowed,
			TimeEnd: time.Now().AddDate(0, periodBuy, 0)}, &db.LinkFilter{Link: key})
		if err != nil {
			tx.Rollback()
			log.Println(op, err)
			return mm.start(user)
		}

		lnks, err := ls.Find(&db.LinkFilter{Link: key}, &db.OrderLinks{}, 0, 1)
		if err != nil {
			log.Println(op, err)
			return mm.start(user)
		}
		if len(lnks) == 0 {
			log.Println(op, errors.New(fmt.Sprintf("Link %s not found", key)))
			return mm.start(user)
		}
		err = mm.vpn.UpdateKey(lnks[0])
		if err != nil {
			log.Println(op, err)
			return mm.start(user)
		}
		user.BalanceAllTime -= priceBuy
		msgText = fmt.Sprintf(MsgBuyOkText,
			mm.vpn.KeyBuilder(lnks[0]),
			time.Now().AddDate(0, periodBuy, 0).Format(layout),
		)
		keys = append(keys, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmHowToUseText, BtmHowToUseValue)),
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmOkText, BtmCancelValue)),
		)
		phs := mm.sm.NewPaymentHistoryStorage()
		tx, _, err = phs.Insert(tx, &models.PaymentHistory{UserId: user.UserId, Amount: priceBuy, TransactionType: BalanceExt})
		if err != nil {
			tx.Rollback()
			log.Println(op, err)
		}
		err = tx.Commit()
		if err != nil {
			log.Println(op, err)
			return mm.start(user)
		}
	}

	msg = tgbotapi.NewMessage(user.UserId, msgText) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	msg.ReplyMarkup = &numericKeyboard1
	return msg
}
