package messages

import (
	"log"
	"vpngigabot/internal/models"
	"vpngigabot/internal/paysystem"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (mm *Messages) chackPay(user *models.User) tgbotapi.Chattable {
	const op = "messages.chackPay"

	user.SetIgnorDeletion(true)
	status, err := mm.ps.CheckPayStatus(&models.Pay{PayId: user.GetCommandData()})
	if err != nil || status == paysystem.YAREFUSED || status == paysystem.YAPROGRESS {
		ans := tgbotapi.NewCallbackWithAlert(user.GetCallBackId(), MsgPayNotFoundText)
		log.Println(op, err)
		// ans := tgbotapi.CallbackConfig{
		// 	CallbackQueryID: user.GetCallBackId(),
		// 	ShowAlert:       true,
		// 	Text:            MsgPayNotFoundText,
		// }
		return ans
	} else if status == paysystem.YASUCCED {
		ans := tgbotapi.CallbackConfig{
			CallbackQueryID: user.GetCallBackId(),
			ShowAlert:       true,
			Text:            MsgPayFoundText,
		}
		return ans
	} else {

		log.Println(op, status)
		ans := tgbotapi.CallbackConfig{
			CallbackQueryID: user.GetCallBackId(),
			ShowAlert:       true,
			Text:            MsgPayNotFoundText,
		}
		return ans
	}

}
