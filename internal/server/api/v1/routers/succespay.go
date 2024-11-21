package routers

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"regexp"
	"vpngigabot/internal/config"
	"vpngigabot/internal/db"
	"vpngigabot/internal/models"
	paysistem "vpngigabot/internal/paysystem"

	"github.com/go-chi/chi"
)

const (
	BalanceAdd = "addmoney"
)

func Succes(cfg *config.Config, sm SotrageManager, de DataExchenge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "routers.Succes"
		log.Println(op, "------------------------Enter------------------------")
		defer func() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Notification received"))
		}()

		iph := sm.NewIpStoriesStorage()
		pls := sm.NewPayLinkStorage()
		us := sm.NewUsersStorage()
		phs := sm.NewPaymentHistoryStorage()

		defer func() {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				//	log.Println(op, "Error:", err)
				ip = r.RemoteAddr
			}
			if err != nil {
				//	log.Println(op, "Error:", err)
				ip = r.RemoteAddr
			}
			txx, _, err := iph.Insert(nil, &models.IpStory{Ip: ip, Link: "postreq", UserAgent: r.UserAgent()})
			if err != nil {
				txx.Rollback()
				log.Println(op, "ip insert", err)
			}
			err = txx.Commit()
			if err != nil {
				log.Println(op, "ip insert", err)
			}
		}()
		if err := r.ParseForm(); err != nil {
			log.Println(op, "Error parsing form:", err)
			return
		}
		lable := ""

		// Выводим все переданные параметры в консоль
		for key, values := range r.Form {
			// fmt.Println(key, values)
			if key == "label" {
				if len(values) >= 1 {
					lable = values[0]
				}
			}
		}
		if lable == "" {
			return
		}

		// Находим платежную ссылку
		pl, err := pls.Find(&db.PayLinkFilter{PayId: lable})
		if err != nil {
			// TODO Pars template error
			log.Println(op, "pls.Find", err)
			return
		}
		// Находим пользователя
		user, err := us.FindOne(pl.UserId)
		if err != nil {
			// TODO Pars template error
			log.Println(op, "us.FindOne", err)
			return
		}
		// Обновляем его баланс
		user.BalanceAllTime += pl.Amount
		tx, err := us.UpdateUser(nil, user, pl.UserId)
		if err != nil {
			// TODO Pars template error
			tx.Rollback()
			log.Println(op, "us.UpdateUser", err)
			return
		}

		// Обновляем историю платежей
		tx, _, err = phs.Insert(tx, &models.PaymentHistory{UserId: user.UserId, Amount: pl.Amount, TransactionType: BalanceAdd})
		if err != nil {
			// TODO Pars template errors
			log.Println(op, "us.UpdateUser", err)
			return

		}

		tx, err = pls.Delete(tx, &db.PayLinkFilter{PayId: lable, Status: paysistem.YAPROGRESS})
		if err != nil {
			// TODO Pars template error
			tx.Rollback()
			log.Println(op, "pls.Delete", err)
			return

		}
		tx.Commit()
		if err != nil {
			log.Println(op, "us.UpdateUser", err)
			return
		}
		log.Println(op, "UserId:", user.UserId, "Success Pay:", pl.Amount)

		datEx := models.DataExchenge{UserId: user.UserId, Command: "VPNKey", Amount: pl.Amount}
		de.Write(&datEx)
		log.Println(op, "******** Success Pay ******** payid:", lable)
	}

}

func SuccesPay(cfg *config.Config, sm SotrageManager, de DataExchenge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "routers.SuccesPay"
		log.Println(op, "!!!!!!!!!!!!!!1!!!!!!!!!!!!!!!!!!!!!!!!")

		pattern := fmt.Sprintf("^[A-Za-z0-9]{%d}$", cfg.LinkLen)
		iph := sm.NewIpStoriesStorage()
		pls := sm.NewPayLinkStorage()
		us := sm.NewUsersStorage()
		phs := sm.NewPaymentHistoryStorage()

		link := chi.URLParam(r, "id")

		defer func() {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				//	log.Println(op, "Error:", err)
				ip = r.RemoteAddr
			}
			if err != nil {
				//	log.Println(op, "Error:", err)
				ip = r.RemoteAddr
			}
			txx, _, err := iph.Insert(nil, &models.IpStory{Ip: ip, Link: link, UserAgent: r.UserAgent()})
			if err != nil {
				txx.Rollback()
				log.Println(op, "ip insert", err)
			}
			err = txx.Commit()
			if err != nil {
				log.Println(op, "ip insert", err)
			}

		}()
		// Валидируем входные данные
		valid := regexp.MustCompile(pattern).MatchString(link)
		if !valid {
			// TODO Pars template error
			log.Println(op, "no valid link")
			return
		}
		// Находим платежную ссылку
		pl, err := pls.Find(&db.PayLinkFilter{PayId: link, Status: paysistem.YAPROGRESS})
		if err != nil {
			// TODO Pars template error
			log.Println(op, "pls.Find", err)
			return
		}
		// Находим пользователя
		user, err := us.FindOne(pl.UserId)
		if err != nil {
			// TODO Pars template error
			log.Println(op, "us.FindOne", err)
			return
		}
		// Обновляем его баланс
		user.BalanceAllTime += pl.Amount
		tx, err := us.UpdateUser(nil, user, pl.UserId)
		if err != nil {
			// TODO Pars template error
			tx.Rollback()
			log.Println(op, "us.UpdateUser", err)
			return
		}

		// Обновляем историю платежей
		tx, _, err = phs.Insert(tx, &models.PaymentHistory{UserId: user.UserId, Amount: pl.Amount, TransactionType: BalanceAdd})
		if err != nil {
			// TODO Pars template errors
			log.Println(op, "us.UpdateUser", err)
			return

		}

		tx, err = pls.Delete(tx, &db.PayLinkFilter{PayId: link, Status: paysistem.YAPROGRESS})
		if err != nil {
			// TODO Pars template error
			tx.Rollback()
			log.Println(op, "pls.Delete", err)
			return

		}
		tx.Commit()
		if err != nil {
			log.Println(op, "us.UpdateUser", err)
			return
		}
		log.Println(op, "UserId:", user.UserId, "Success Pay:", pl.Amount)
		// TODO Pars template error
		res := struct {
			Amount      int
			TotalAmount int
		}{Amount: pl.Amount, TotalAmount: user.BalanceAllTime}
		parsedTemplate, err := template.ParseFiles("template/successpay.html")
		if err != nil {
			log.Println(op, err)
			return
		}
		err = parsedTemplate.Execute(w, res)
		if err != nil {
			log.Println("Error executing template :", err)
			return
		}
		datEx := models.DataExchenge{UserId: user.UserId, Command: "VPNKey", Amount: res.Amount}
		de.Write(&datEx)
	}

}
