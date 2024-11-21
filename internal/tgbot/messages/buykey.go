package messages

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	BalanceExt string = "BalanceExt"
)

func (mm *Messages) buyKey(user *models.User) tgbotapi.MessageConfig {
	const op = "messages.priceMeny"
	var msg tgbotapi.MessageConfig
	msgText := ""
	var keys [][]tgbotapi.InlineKeyboardButton

	cd := strings.Split(user.GetCommandData(), "-")
	if len(cd) != 2 {
		log.Println(op, errors.New("Len command data less 2 "))
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
	if user.BalanceAllTime < priceBuy {
		msgText = MsgBuyErrText
		keys = append(keys,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(BtmOkText, BtmAddKeyValue)),
		)
	} else {
		key, err := mm.vpn.CreateKey(user)
		if err != nil {
			log.Println(op, err)
			return mm.start(user)
		}
		// key.Link = randString(mm.cfg.LinkLen)
		key.State = models.StateAllowed
		key.TimeEnd = time.Now().AddDate(0, periodBuy, 0)
		ls := mm.sm.NewLinkStorage()
		tx, _, err := ls.Insert(nil, key)
		if err != nil {
			tx.Rollback()
			log.Println(op, err)
			return mm.start(user)
		}
		user.BalanceAllTime -= priceBuy
		msgText = fmt.Sprintf(MsgBuyOkText,
			mm.vpn.KeyBuilder(key),
			key.TimeEnd.Format(layout),
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
		}
	}

	msg = tgbotapi.NewMessage(user.UserId, msgText) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	msg.ReplyMarkup = &numericKeyboard1
	return msg
}
