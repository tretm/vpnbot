package messages

import (
	"fmt"
	"html"
	"log"
	"strconv"
	"vpngigabot/internal/db"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) keyMenu(user *models.User) tgbotapi.MessageConfig {
	const op = "messages.keyMenu"
	var keys [][]tgbotapi.InlineKeyboardButton
	ls := mm.sm.NewLinkStorage()
	lnk, err := ls.Find(&db.LinkFilter{UserId: strconv.FormatInt(user.UserId, 10)}, &db.OrderLinks{Date: true}, 0, mm.cfg.RowsLimit)
	if err != nil {
		log.Println(op, err)
	} else {
		for i, l := range lnk {
			state := ""
			if l.State == models.StateAllowed {
				state = html.UnescapeString("&#9989;")
			} else {
				state = html.UnescapeString("&#10060;")
			}
			keys = append(keys, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf(BtmVpnKeyDetaleText, state, i+1, l.Link),
					BtmVpnKeyDetaleValue+"_"+l.Link)))
		}
	}

	keys = append(keys, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmAddKeyText, BtmAddKeyValue)),

		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmCancelText, BtmCancelValue)),
	)
	msg := tgbotapi.NewMessage(user.UserId, fmt.Sprintf(MsgKeysMenu, user.BalanceAllTime)) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
	numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(keys...)
	msg.ReplyMarkup = &numericKeyboard1
	return msg
}
