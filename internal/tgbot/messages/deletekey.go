package messages

import (
	"log"
	"strconv"
	"vpngigabot/internal/db"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) deleteKey(user *models.User) *tgbotapi.MessageConfig {
	const op = "deleteKey"
	textMsg := ""
	err := mm.vpn.DeleteKey(user.GetCommandData())
	if err != nil {
		textMsg = "Ошибка"
		log.Println(op, err)
	} else {
		textMsg = "Ключ удален"
		log.Println(op, "key: removed")
	}
	ls := mm.sm.NewLinkStorage()
	lnk := models.Link{}
	lnk.State = models.StateBanned
	uId := strconv.FormatInt(user.UserId, 10)
	tx, _, err := ls.Update(nil, &lnk, &db.LinkFilter{UserId: uId})
	if err != nil {
		tx.Rollback()
		log.Println(op, err)
	}
	err = tx.Commit()
	if err != nil {
		log.Println(op, err)
	}

	// bot.DeleteMessage(tgbotapi.NewDeleteMessage(p.user.Userid, p.user.Messageid))
	msg := tgbotapi.NewMessage(user.UserId, textMsg) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)
	// numericKeyboard1 := tgbotapi.NewInlineKeyboardMarkup(
	// 	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmVpnKeyText, BtmVpnKeyValue)),
	// 	// tgbotapi.NewInlineKeyboardButtonData(filds.BtmExitmoney, filds.ExitMoney)),
	// 	// tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(filds.BtmRuue, filds.Ruue),
	// 	// 	tgbotapi.NewInlineKeyboardButtonData(filds.BtmUeru, filds.Ueru)),
	// 	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmBalancAddText, BtmBalancAddValue)),
	// 	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmHistoryPayText, BtmHistoryPayValue)),
	// 	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmHelpText, BtmHelpValue)),
	// 	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BtmReferalText, BtmReferalValue)),
	// )
	// msg.ReplyMarkup = &numericKeyboard1
	return &msg
}
