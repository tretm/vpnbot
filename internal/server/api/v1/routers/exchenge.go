package routers

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"vpngigabot/internal/config"
	"vpngigabot/internal/models"

	"github.com/go-chi/render"
)

type DataExchenge interface {
	Write(*models.DataExchenge)
	Read() chan models.DataExchenge
}

type Resp struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `jsong:"message"`
}

func SendMsgToBot(cfg *config.Config, sm SotrageManager, de DataExchenge) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "routers.SendMsgToBot"

		iph := sm.NewIpStoriesStorage()
		defer func() {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				//	log.Println(op, "Error:", err)
				ip = r.RemoteAddr
			}
			txx, _, err := iph.Insert(nil, &models.IpStory{Ip: ip, Link: "/ex", UserAgent: r.UserAgent()})
			if err != nil {
				txx.Rollback()
				log.Println(op, "ip insert", err)
			}
			err = txx.Commit()
			if err != nil {
				log.Println(op, "ip insert", err)
			}
		}()
		datEx := models.DataExchenge{}
		if err := json.NewDecoder(r.Body).Decode(&datEx); err != nil {
			render.JSON(w, r, &Resp{Code: 500, Status: "error", Message: err.Error()})
			log.Println(op, "NewDecoder", err)
			return
		}
		defer func() { _ = r.Body.Close() }()

		de.Write(&datEx)
		render.JSON(w, r, &Resp{Code: 200, Status: "ok"})

	}

}
