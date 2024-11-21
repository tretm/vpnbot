package messages

import (
	"strings"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) confirmExtBuy(user *models.User) tgbotapi.MessageConfig {
	const op = "messages.confirmExtBuy"
	c := strings.Count(user.GetCommandData(), "-")
	if c == 1 {
		return mm.buyKey(user)
	} else if c == 2 {
		return mm.buyExtendKey(user)
	} else {
		return mm.start(user)
	}
}
