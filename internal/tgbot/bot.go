package bot

import (
	"fmt"
	"log"
	"vpngigabot/internal/config"
	"vpngigabot/internal/db"
	"vpngigabot/internal/models"

	// tgbotapi "github.com/Syfaro/telegram-bot-api"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//"333308760:AAGkA6o4ldOxw1y_ioSGtP8cUanYbFKTMaM") //inidat.BotToken) //"333308760:AAGkA6o4ldOxw1y_ioSGtP8cUanYbFKTMaM")

type TgBot struct {
	bot    *tgbotapi.BotAPI
	update *tgbotapi.UpdatesChannel
	cfg    *config.Config
}
type VPNManager interface {
	CreateKey(*models.User) (*models.Link, error)
	SetTraficLimit(key string, limitGb int) error
	DeleteKey(string) error
	UpdateKey(*models.Link) error
	KeyBuilder(*models.Link) string
}

func NewTgBot(iniDat *config.Config) (*TgBot, error) {
	var tgBot TgBot
	bot, err := tgbotapi.NewBotAPI(iniDat.BotToken)
	if err != nil {
		return nil, err
	}
	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	upd := bot.GetUpdatesChan(ucfg)
	tgBot.bot = bot
	tgBot.update = &upd
	tgBot.cfg = iniDat
	return &tgBot, err
}

type SotrageManager interface {
	NewLinkStorage() db.Links
	NewUsersStorage() db.Users
	NewPaymentHistoryStorage() db.PaymentHistory
	NewPayLinkStorage() db.PayLink
}

type TgMethods interface {
	MessageManager(user *models.User) tgbotapi.Chattable
}
type DataExchenge interface {
	Write(*models.DataExchenge)
	Read() chan models.DataExchenge
}

func (tgBot *TgBot) RunBot(sm SotrageManager, vpn VPNManager, tm TgMethods, de DataExchenge) {
	const op = "bot.RunBot"
	us := sm.NewUsersStorage()
	amount := 0
	tgBot.scheduler(sm, vpn, tm)
	go func() {
		// if tgBot.cfg.SendNotify != 0 {
		// 	go sendNotify(tgBot.cfg.AdminId, tgBot.bot, us)
		// }
		for {
			userforsend := models.NewUser()
			update := tgbotapi.Update{}
			select {
			case update = <-*tgBot.update:

				// Получаем тип сообщения Callback или обычное Message устанавливаем идентификатор и логин в модель пользователя

				if update.CallbackQuery != nil {
					userforsend.SetUserId(update.CallbackQuery.Message.Chat.ID).SetUserName(getlogin(update.CallbackQuery.Message.Chat))
				} else if update.Message != nil {
					userforsend.SetUserId(update.Message.Chat.ID).SetUserName(getlogin(update.Message.Chat))
				} else if update.PreCheckoutQuery != nil {
					userforsend.SetUserId(int64(update.PreCheckoutQuery.From.ID))
				}

			case datEx := <-de.Read():
				userforsend.SetUserId(datEx.UserId)
				amount = datEx.Amount
			}
			if userforsend.UserId <= 0 {
				continue
			}

			go func() {

				// Ищем пользователя по его Идентификатору, если не находим то производим вставку в базу данных (создаем нового пользователя)
				user, err := us.FindOne(userforsend.UserId)
				if err != nil {
					log.Println(op, "FindOne", err)
					tx, _, err := us.Insert(nil, userforsend)
					if err != nil {
						tx.Rollback()
						log.Println(op, "Insert", err)
						return
					}
					tx.Commit()
					if err != nil {
						log.Println(op, "Insert", err)
						return
					}
					//TODO command in User
					// userforsend.Command=
					// tgBot.sendMsgToUser(userforsend, tm.MessageManager(userforsend))
					tgBot.updateManager(userforsend, update, sm, tm)
					return

				}
				// Если пользователь заблокирован отправляем сообщение в канал и игнорируем пользователя так как будто бот не работает
				if user.Status == models.BANNED {
					err = tgBot.sendLogToChan(tgBot.cfg.ReportChanId, "Banned user", user)
					if err != nil {
						log.Println(op, "sendLogTochan", err)
					}
					return
				}

				if user.UserName != userforsend.UserName {
					user.HistoryUserName += "," + user.UserName
					user.UserName = userforsend.UserName

				}
				userforsend = user
				// выбираем что делать в зависимости от того что в update
				userforsend.Balance = amount
				tgBot.updateManager(userforsend, update, sm, tm)
			}()
		}

	}()
}

// func (tgBot *TgBot) RunBot(sm SotrageManager, vpn VPNManager, tm TgMethods, de DataExchenge) {
// 	const op = "bot.RunBot"
// 	us := sm.NewUsersStorage()
// 	tgBot.scheduler(sm, vpn, tm)
// 	go func() {
// 		// if tgBot.cfg.SendNotify != 0 {
// 		// 	go sendNotify(tgBot.cfg.AdminId, tgBot.bot, us)
// 		// }
// 		for update := range *tgBot.update {

// 			// Получаем тип сообщения Callback или обычное Message устанавливаем идентификатор и логин в модель пользователя
// 			userforsend := models.NewUser()
// 			if update.CallbackQuery != nil {
// 				userforsend.SetUserId(update.CallbackQuery.Message.Chat.ID).SetUserName(getlogin(update.CallbackQuery.Message.Chat))
// 			} else if update.Message != nil {
// 				userforsend.SetUserId(update.Message.Chat.ID).SetUserName(getlogin(update.Message.Chat))
// 			} else if update.PreCheckoutQuery != nil {
// 				userforsend.SetUserId(int64(update.PreCheckoutQuery.From.ID))
// 			}
// 			if userforsend.UserId <= 0 {
// 				continue
// 			}
// 			// Ищем пользователя по его Идентификатору, если не находим то производим вставку в базу данных (создаем нового пользователя)
// 			user, err := us.FindOne(userforsend.UserId)
// 			if err != nil {
// 				log.Println(op, "FindOne", err)
// 				tx, _, err := us.Insert(nil, userforsend)
// 				if err != nil {
// 					tx.Rollback()
// 					log.Println(op, "Insert", err)
// 					continue
// 				}
// 				tx.Commit()
// 				if err != nil {
// 					log.Println(op, "Insert", err)
// 					continue
// 				}
// 				//TODO command in User
// 				// userforsend.Command=
// 				tgBot.sendMsgToUser(userforsend, tm.MessageManager(userforsend))
// 				continue

// 			}

// 			// Если пользователь заблокирован отправляем сообщение в канал и игнорируем пользователя так как будто бот не работает
// 			if user.Status == models.BANNED {
// 				err = tgBot.sendLogToChan(tgBot.cfg.ReportChanId, "Banned user", user)
// 				if err != nil {
// 					log.Println(op, "sendLogTochan", err)
// 				}
// 				continue
// 			}

// 			if user.UserName != userforsend.UserName {
// 				user.HistoryUserName += "," + user.UserName
// 				user.UserName = userforsend.UserName

// 			}
// 			userforsend = user
// 			// выбираем что делать в зависимости от того что в update
// 			tgBot.updateManager(userforsend, update, sm, tm)

// 		}
// 	}()
// }

func (tgBot *TgBot) RunExperementBot(sm SotrageManager, tm TgMethods) {
	const op = "bot.RunExperementBot"

	go func() {
		// if tgBot.cfg.SendNotify != 0 {
		// 	go sendNotify(tgBot.cfg.AdminId, tgBot.bot, us)
		// }
		for update := range *tgBot.update {

			// Получаем тип сообщения Callback или обычное Message устанавливаем идентификатор и логин в модель пользователя
			userforsend := models.NewUser()
			if update.CallbackQuery != nil {
				userforsend.SetUserId(update.CallbackQuery.Message.Chat.ID).SetUserName(getlogin(update.CallbackQuery.Message.Chat))
			} else if update.Message != nil {
				userforsend.SetUserId(update.Message.Chat.ID).SetUserName(getlogin(update.Message.Chat))
				if update.Message.SuccessfulPayment != nil {
					fmt.Println("+++++++++++++", update.Message.SuccessfulPayment)
				}

			} else if update.PreCheckoutQuery != nil {
				log.Println("*/-/", update.PreCheckoutQuery.From.ID)
				pca := tgbotapi.PreCheckoutConfig{
					OK:                 true,
					PreCheckoutQueryID: update.PreCheckoutQuery.ID,
				}
				_, err := tgBot.bot.Send(pca)
				if err != nil {
					log.Println("+++", op, err)
				}

			}

			if userforsend.UserId <= 0 {
				continue
			}
			if update.Message != nil {
				switch update.Message.Command() {
				case "start":
					// create invoice
					invoice := tgbotapi.NewInvoice(userforsend.UserId, "Test invoice", "description here", "custom_payload",
						"381764678:TEST:73448", "start_param", "RUB",
						[]tgbotapi.LabeledPrice{{Label: "RUB", Amount: 200000}})
					//"1744374395:TEST:bdbc83a1d72e7b91c235"

					_, err := tgBot.bot.Send(invoice)
					if err != nil {
						log.Println("!!!", op, err)
					}

				}
			}

		}
	}()
}
