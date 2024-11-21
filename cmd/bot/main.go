package main

import (
	"fmt"
	"log"
	"vpngigabot/internal/config"
	"vpngigabot/internal/dataexchenge"
	"vpngigabot/internal/db/mysql"
	"vpngigabot/internal/paysystem/youmoney"
	"vpngigabot/internal/server"
	v1 "vpngigabot/internal/server/api/v1"
	tgbot "vpngigabot/internal/tgbot"
	"vpngigabot/internal/tgbot/messages"
	"vpngigabot/internal/vpn"
)

func main() {

	cfg, err := config.PrepareConfig()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("VpnGigaBot")
	msdb := &mysql.Manager{}
	msdb.Connect(cfg.StorageConfig)
	//vpnMng := vpn.New("outline", cfg)
	vpnMng := vpn.New(vpn.VtoRay, cfg)
	tb, err := tgbot.NewTgBot(cfg)
	if err != nil {
		log.Println(err)
		return
	}
	ym, err := youmoney.New(cfg)
	if err != nil {
		log.Println(err)
		return
	}

	mm := messages.New(cfg, vpnMng, msdb, ym)
	de := dataexchenge.New()
	//tb.RunExperementBot(msdb, mm)
	tb.RunBot(msdb, vpnMng, mm, de)
	rout := v1.CreateRouter(cfg, msdb, nil, de)
	fmt.Println("Is runned")
	srv := server.NewServer(cfg, rout)
	srv.Serve()

	fmt.Println("Is runned")

}
