package bot

import (
	"fmt"
	"log"
	"strings"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	layout    = "2006-01-02"
	SuccesPay = "SuccesPay"
	PreCheck  = "PreCheck"
)

func getlogin(upd *tgbotapi.Chat) (login string) {
	if upd.UserName != "" {
		return upd.UserName
	}
	// if upd.FirstName == "" {
	// 	return upd.LastName
	// }
	// if upd.LastName == "" {
	// 	return upd.FirstName
	// }
	return strings.TrimSpace(upd.FirstName + " " + upd.LastName)
}

func (tgBot *TgBot) updateManager(userForSend *models.User, update tgbotapi.Update, sm SotrageManager, tm TgMethods) {
	const op = "bot.RunBot"
	us := sm.NewUsersStorage()
	switch {
	case update.PreCheckoutQuery != nil:
		pca := tgbotapi.PreCheckoutConfig{
			OK:                 true,
			PreCheckoutQueryID: update.PreCheckoutQuery.ID,
		}
		_, err := tgBot.bot.Send(pca)
		if err != nil {
			log.Println(op, err)
		}
	case update.CallbackQuery != nil:
		go func() {
			// Делаяется для того чтобы перестала моргать нажатая кнопка
			// _, err := tgBot.bot.Send(tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID})
			// if err != nil {
			// 	log.Println(op, "AnswerCallbackQuery", err)
			// }
			if command := update.CallbackQuery.Data; command != "" {
				userForSend.SetCommand(command).SetCallBackId(update.CallbackQuery.ID)
				tgBot.sendMsgToUser(userForSend, tm.MessageManager(userForSend))
				tx, err := us.UpdateUser(nil, userForSend, userForSend.UserId)
				if err != nil {
					tx.Rollback()
					log.Println(op, "UpdateUser", err)
				}
				err = tx.Commit()
				if err != nil {
					log.Println(op, "UpdateUser", err)
				}
				err = tgBot.sendLogToChan(tgBot.cfg.ReportChanId, command, userForSend)
				if err != nil {
					log.Println("bot.sendLogToChan", err)

				}
			}
		}()
	case update.Message != nil:
		switch {
		case update.Message.Command() != "":
			go func() {
				userForSend.SetCommand(update.Message.Command())
				tgBot.sendMsgToUser(userForSend, tm.MessageManager(userForSend))
				tx, err := us.UpdateUser(nil, userForSend, userForSend.UserId)
				if err != nil {
					tx.Rollback()
					log.Println(op, "UpdateUser", err)
				}
				tx.Commit()
				if err != nil {
					log.Println(op, "UpdateUser", err)
				}
				err = tgBot.sendLogToChan(tgBot.cfg.ReportChanId, userForSend.GetCommand(), userForSend)
				if err != nil {
					log.Println("bot.sendLogToChan", err)
				}
			}()
		case update.Message.SuccessfulPayment != nil:
			go func() {
				userForSend.SetCommand(SuccesPay)
				userForSend.Balance = update.Message.SuccessfulPayment.TotalAmount / 100
				userForSend.BalanceAllTime += userForSend.Balance
				tgBot.sendMsgToUser(userForSend, tm.MessageManager(userForSend))
				tx, err := us.UpdateUser(nil, userForSend, userForSend.UserId)
				if err != nil {
					tx.Rollback()
					log.Println(op, "UpdateUser", err)
				}
				tx.Commit()
				if err != nil {
					log.Println(op, "UpdateUser", err)
				}
				err = tgBot.sendLogToChan(tgBot.cfg.ReportChanId, userForSend.GetCommand(), userForSend)
				if err != nil {
					log.Println("bot.sendLogToChan", err)
				}
			}()

		default:
			go func() {
				// TODO command for user
				tgBot.sendMsgToUser(userForSend, tm.MessageManager(userForSend))
				tx, err := us.UpdateUser(nil, userForSend, userForSend.UserId)
				if err != nil {
					tx.Rollback()
					log.Println(op, "UpdateUser", err)
				}
				tx.Commit()
				if err != nil {
					log.Println(op, "UpdateUser", err)
				}
				err = tgBot.sendLogToChan(tgBot.cfg.ReportChanId, update.Message.Text, userForSend)
				if err != nil {
					log.Println("bot.sendLogToChan", err)
				}
				log.Println("User:", userForSend.UserId, "Send: SOMEMESSAGE")
			}()
		}
	default:

		go func() {
			// TODO command for user
			userForSend.SetCommand("start")
			tgBot.sendMsgToUser(userForSend, tm.MessageManager(userForSend))
			tx, err := us.UpdateUser(nil, userForSend, userForSend.UserId)
			if err != nil {
				tx.Rollback()
				log.Println(op, "UpdateUser", err)
			}
			tx.Commit()
			if err != nil {
				log.Println(op, "UpdateUser", err)
			}
			err = tgBot.sendLogToChan(tgBot.cfg.ReportChanId, fmt.Sprintf("Подтверждение платежа: %d", userForSend.Balance), userForSend)
			if err != nil {
				log.Println("bot.sendLogToChan", err)
			}
			// log.Println("User:", userForSend.UserId, "Send: Pay confirmed")
		}()
	}
}
