package messages

import (
	"fmt"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) notification(user *models.User) tgbotapi.MessageConfig {

	msg := tgbotapi.NewMessage(user.UserId, fmt.Sprintf(MsgNotificationText,
		user.GetCommandData())) //NewEditMessage(p.user.Userid, m.MessageID, numericKeyboard)

	return msg
}
