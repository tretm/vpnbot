package routers

import (
	"log"
	"net"
	"net/http"
	"vpngigabot/internal/config"
	"vpngigabot/internal/db"
	"vpngigabot/internal/models"

	"github.com/go-chi/render"
)

type SotrageManager interface {
	NewLinkStorage() db.Links
	NewUsersStorage() db.Users
	NewPaymentHistoryStorage() db.PaymentHistory
	NewIpStoriesStorage() db.IpStories
	NewPayLinkStorage() db.PayLink
}

type IndexOk struct {
	Version      string `json:"v"`
	ServerStatus string `json:"status"`
}

func Index(cfg *config.Config, sm SotrageManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "routers.Index"
		iph := sm.NewIpStoriesStorage()
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
			tx, _, err := iph.Insert(nil, &models.IpStory{Ip: ip, Link: "/", UserAgent: r.UserAgent()})
			if err != nil {
				tx.Rollback()
				log.Println(op, "ip insert", err)
			}
			err = tx.Commit()
			if err != nil {
				log.Println(op, "ip insert", err)
			}

		}()
		w.Header().Set("Content-Type", "application/json")
		render.JSON(w, r, IndexOk{Version: "1.0.0", ServerStatus: "OK"})
	}

}
