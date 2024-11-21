package messages

import (
	"fmt"
	"log"
	"time"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) newOutlineKey(user *models.User) tgbotapi.MessageConfig {
	const op = "newOutlineKey"
	textMsg := ""
	var keys [][]tgbotapi.InlineKeyboardButton
	keys = append(keys, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmHowToUseText, BtmHowToUseValue)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmOkText, BtmCancelValue)),
	)
	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	ls := mm.sm.NewLinkStorage()

	key, err := mm.vpn.CreateKey(user)
	if err != nil {
		log.Println(op, err)
		return mm.start(user)
	}
	// key.Link = randString(mm.cfg.LinkLen)
	key.State = models.StateAllowed
	key.TimeEnd = time.Now().AddDate(0, 0, mm.cfg.TrialPeriod)
	tx, _, err := ls.Insert(nil, key)
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
	user.TestUsed = true

	textMsg = fmt.Sprintf(MessageTrialGood,
		mm.cfg.TrialPeriod,
		mm.vpn.KeyBuilder(key))

	msg := tgbotapi.NewMessage(user.UserId, textMsg) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
	msg.ReplyMarkup = &numericKeyboard1
	return msg
}
