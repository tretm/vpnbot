package messages

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"vpngigabot/internal/db"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) keyDetalisation(user *models.User) tgbotapi.MessageConfig {
	const op = "newOutlineKey"
	textMsg := ""
	var keys [][]tgbotapi.InlineKeyboardButton

	ls := mm.sm.NewLinkStorage()
	// проверяем нетли уже существующего ключа у пользователя
	avaliblKeys, err := ls.Find(&db.LinkFilter{UserId: strconv.FormatInt(user.UserId, 10), Link: user.GetCommandData()}, &db.OrderLinks{}, 0, 1)
	if err != nil && err.Error() != "links not found" {
		log.Println(op, err)
		return mm.start(user)
	}
	if len(avaliblKeys) > 0 {
		if avaliblKeys[0].State == models.StateAllowed {
			textMsg = fmt.Sprintf(MsgKeyDetalisationOk,
				mm.vpn.KeyBuilder(avaliblKeys[0]), avaliblKeys[0].TimeEnd.Format(layout),
				int(avaliblKeys[0].TimeEnd.Sub(time.Now()).Hours()/24)+1)
		} else {
			textMsg = fmt.Sprintf(MsgKeyDetalisationBad,
				mm.vpn.KeyBuilder(avaliblKeys[0]), avaliblKeys[0].TimeEnd.Format(layout),
			)
			keys = append(keys,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(BtmExtendText, BtmExtndValue+"_"+avaliblKeys[0].Link)))
		}

	}
	keys = append(keys, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmHowToUseText, BtmHowToUseValue)),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmOkText, BtmVpnKeyValue)),
	)
	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	msg := tgbotapi.NewMessage(user.UserId, textMsg) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
	msg.ReplyMarkup = &numericKeyboard1
	return msg
}
