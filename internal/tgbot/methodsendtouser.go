package bot

import (
	"log"
	"time"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tgBot *TgBot) sendMsgToUser(user *models.User, msg tgbotapi.Chattable) {
	switch msg.(type) {
	case tgbotapi.MessageConfig:
		msg := msg.(tgbotapi.MessageConfig)
		msg.ParseMode = "html" //"markdown"
		msg.DisableWebPagePreview = true

		tgBot.bot.Send(tgbotapi.NewDeleteMessage(user.UserId, int(user.LastMsgId)))

		resMsg, err := tgBot.bot.Send(msg)
		if err != nil {
			log.Println("sendMsgToUser", err)
			time.Sleep(1 * time.Second)
			resMsg, err = tgBot.bot.Send(msg)
			if err != nil {
				log.Println("sendMsgToUser 2", err)
			} else {
				log.Println("sendMsgToUser Message sended to user msg id", resMsg.MessageID)
			}

		}

		user.SetLastMsgId(int64(resMsg.MessageID))

	case tgbotapi.InvoiceConfig:
		_, err := tgBot.bot.Send(msg)
		if err != nil {
			log.Println("sendMsgToUser", err)
		}
	case tgbotapi.CallbackConfig:
		msg := msg.(tgbotapi.CallbackConfig)

		_, err := tgBot.bot.Send(msg)
		if err != nil {
			log.Println("sendMsgToUser", err)
		}

	}

}
func (tgBot *TgBot) sendPhotoMsgToUser(userId int64, msgId int64, message, btmCode string) {
	// msg := tgbotapi.NewMessage(userId, message)
	// //fmt.Println(message)
	// var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	// 	tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(BOTTOMDESCRIPTION, btmCode)),
	// )
	// msg.ReplyMarkup = numericKeyboard
	// msg.ParseMode = "markdown"
	// msg.ReplyToMessageID = int(msgId)
	// _, err := tgBot.bot.Send(msg)
	// if err != nil {
	// 	log.Println("sendMsgToUser", userId, err)
	// }
}
