package messages

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) newTrialOutlineKey(user *models.User) tgbotapi.MessageConfig {
	const op = "newOutlineKey"
	textMsg := ""
	var keys [][]tgbotapi.InlineKeyboardButton
	keys = append(keys, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmHowToUseText, BtmHowToUseValue)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmOkText, BtmCancelValue)),
	)
	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	ls := mm.sm.NewLinkStorage()
	// проверяем нетли уже существующего ключа у пользователя
	// avaliblKeys, err := ls.Find(&db.LinkFilter{}, &db.OrderLinks{}, 0, 1)
	// if err != nil && err.Error() != "links not found" {
	// 	log.Println(op, err)
	// 	return mm.start(user)
	// }
	// if len(avaliblKeys) > 0 {
	// 	if avaliblKeys[0].State == models.StateBanned {
	// 		textMsg = fmt.Sprintf(MessageTrialExist,
	// 			mm.cfg.Ssconf+avaliblKeys[0].Link, MessageTrialIsFinished)
	// 	} else {
	// 		textMsg = fmt.Sprintf(MessageTrialExist,
	// 			mm.cfg.Ssconf+avaliblKeys[0].Link, "")
	// 	}
	// 	msg := tgbotapi.NewMessage(user.UserId, textMsg)
	// 	msg.ReplyMarkup = &numericKeyboard1
	// 	return msg
	// }
	// если ключа нет, создаем новый ключ
	key, err := mm.vpn.CreateKey(user)
	if err != nil {
		log.Println("!-!-!-!-", op, err)
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

// func (mm *Messages) newOutlineKey(user *models.User) *tgbotapi.MessageConfig {
// 	const op = "newOutlineKey"
// 	textMsg := ""

// 	sm:=mm.sm.NewLinkStorage()
// 	// проверяем нетли уже существующего ключа у пользователя
// 	sm.Find()

// 	key, err := mm.vpn.CreateKey(user)
// 	if err != nil {
// 		textMsg = "Ошибка"
// 		log.Println(op, err)
// 	}

// 	key.Link = randString(mm.cfg.LinkLen)
// 	key.State = models.StateAllowed

// 	daysStr := user.GetCommandData()
// 	daysInt, err := strconv.Atoi(daysStr)
// 	if err != nil {
// 		log.Println(op, err)
// 		daysInt = mm.cfg.DefaultDays
// 	}

//		key.TimeEnd = time.Now().AddDate(0, 0, daysInt)
//		ls := mm.sm.NewLinkStorage()
//		_, err = ls.Insert(key)
//		if err != nil {
//			log.Println(op, err)
//		}
//		textMsg = mm.cfg.Ssconf + key.Link
//		// bot.DeleteMessage(tgbotapi.NewDeleteMessage(p.user.Userid, p.user.Messageid))
//		msg := tgbotapi.NewMessage(user.UserId, textMsg) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
//		// numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(
//		// 	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmVpnKeyText, BtmVpnKeyValue)),
//		// 	// tgbotapi.NewInlineKeyboardButtonData(filds.BtmExitmoney, filds.ExitMoney)),
//		// 	// tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(filds.BtmRuue, filds.Ruue),
//		// 	// 	tgbotapi.NewInlineKeyboardButtonData(filds.BtmUeru, filds.Ueru)),
//		// 	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmBalancAddText, BtmBalancAddValue)),
//		// 	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmHistoryPayText, BtmHistoryPayValue)),
//		// 	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmHelpText, BtmHelpValue)),
//		// 	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmReferalText, BtmReferalValue)),
//		// )
//		// msg.ReplyMarkup = &numericKeyboard1
//		return &msg
//	}
func randString(length int) string {

	rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	return str
}
