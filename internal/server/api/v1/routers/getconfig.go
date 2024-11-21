package routers

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"vpngigabot/internal/config"
	"vpngigabot/internal/db"
	"vpngigabot/internal/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// {"server":"keysvout.tech","server_port":"39012","password":"WBH4lag4idAqzjo8xLomDK","method":"chacha20-ietf-poly1305"}

type OutlineConfig struct {
	Server     string `json:"server"`
	ServerPort string `jsong:"server_port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
}
type Vpn interface {
	CreateKey(*models.User) (*models.Link, error)
	DeleteKey(string) error
}

func GetConfig(cfg *config.Config, sm SotrageManager, vpn Vpn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op string = "routers.GetConfig"
		pattern := fmt.Sprintf("^[A-Za-z0-9]{%d}$", cfg.LinkLen)
		iph := sm.NewIpStoriesStorage()
		ls := sm.NewLinkStorage()
		link := chi.URLParam(r, "id")

		defer func() {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				//	log.Println(op, "Error:", err)
				ip = r.RemoteAddr
			}
			tx, _, err := iph.Insert(nil, &models.IpStory{Ip: ip, Link: link, UserAgent: r.UserAgent()})
			if err != nil {
				tx.Rollback()
				log.Println(op, "ip insert", err)
			}
			err = tx.Commit()
			if err != nil {
				log.Println(op, "ip insert", err)
			}

		}()

		valid := regexp.MustCompile(pattern).MatchString(link)
		if !valid {
			render.JSON(w, r, OutlineConfig{})
			log.Println(op, "no valid link")
			return
		}

		conf, err := ls.Find(&db.LinkFilter{Link: link, State: models.StateAllowed}, &db.OrderLinks{}, 0, 1)
		if err != nil {
			render.JSON(w, r, OutlineConfig{})
			log.Println(op, "ls.Find", err)
			return
		}

		if len(conf) > 0 {
			err := vpn.DeleteKey(conf[0].VpnPassword)
			if err != nil {
				log.Println(op, err)

			}
			lnk, err := vpn.CreateKey(models.NewUser().SetUserId(conf[0].UserId))
			if err != nil {
				log.Println(op, err)
				return
			}
			render.JSON(w, r,
				OutlineConfig{Server: cfg.OutlineVpn.DomeName,
					ServerPort: cfg.OutlineVpn.Port,
					Password:   lnk.VpnPassword,
					Method:     cfg.Method})
			return
		}
		render.JSON(w, r, OutlineConfig{})
		log.Println(op, "outline config not found")

	}

}
