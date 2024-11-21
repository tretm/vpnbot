package bot

import (
	"fmt"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tgBot *TgBot) sendLogToChan(idChan int64, message string, usr *models.User) error {

	mess := fmt.Sprintf("Логин\n<code>%s</code>\nId пользователя\n<code>%d</code>\nЗагрузил:\n %s", usr.UserName, usr.UserId, messageToChan(message))

	if idChan != 0 {
		msg := tgbotapi.NewMessage(idChan, mess) //-1001402390028, ) //601421875 сер
		msg.ParseMode = "html"
		_, err := tgBot.bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func messageToChan(message string) string {
	switch message {
	case "start":
		return "Нажал start"
	default:
		return message
	}
}
