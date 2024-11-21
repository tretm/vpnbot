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

func (mm *Messages) confirmBuy(user *models.User) tgbotapi.MessageConfig {
	const op = "messages.confirmBuy"
	var msg tgbotapi.MessageConfig
	msgText := ""

	var keys [][]tgbotapi.InlineKeyboardButton
	cd := strings.Split(user.GetCommandData(), "-")

	if len(cd) != 2 && len(cd) != 3 {
		log.Println(op, errors.New("Len command data less 3 or 2"))
		return mm.start(user)
	} else if len(cd) == 2 {
		period, err := strconv.Atoi(cd[0])
		if err != nil {
			log.Println(op, err)
			return mm.start(user)
		}
		msgText = fmt.Sprintf(MsgConfirmBuyText,
			int(time.Now().AddDate(0, period, 0).Sub(time.Now()).Hours()/24)+1,
			cd[1])

		keys = append(keys, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BtmConfirmYesText, BtmConfirmYesValue+"_"+user.GetCommandData()),
			tgbotapi.NewInlineKeyboardButtonData(BtmConfirmNoText, BtmAddKeyValue),
		),
		)

	} else if len(cd) == 3 {
		period, err := strconv.Atoi(cd[0])
		if err != nil {
			log.Println(op, err)
			return mm.start(user)
		}
		msgText = fmt.Sprintf(MsgConfirmExtendText, cd[2],
			int(time.Now().AddDate(0, period, 0).Sub(time.Now()).Hours()/24)+1,
			cd[1])
		keys = append(keys, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(BtmConfirmYesText, BtmConfirmYesValue+"_"+user.GetCommandData()),
			tgbotapi.NewInlineKeyboardButtonData(BtmConfirmNoText, BtmExtndValue+"_"+cd[2]),
		),
		)
	}

	msg = tgbotapi.NewMessage(user.UserId, msgText) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	msg.ReplyMarkup = &numericKeyboard1
	return msg
}
