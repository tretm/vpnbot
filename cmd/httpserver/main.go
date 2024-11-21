package main

import (
	"log"
	"vpngigabot/internal/config"
	"vpngigabot/internal/dataexchenge"
	"vpngigabot/internal/db/mysql"
	"vpngigabot/internal/server"
	v1 "vpngigabot/internal/server/api/v1"
)

func main() {
	cfg, err := config.PrepareConfig()
	if err != nil {
		log.Println(err)
		return
	}

	msdb := &mysql.Manager{}
	msdb.Connect(cfg.StorageConfig)
	// outlinevpn := outline.New(cfg)
	de := dataexchenge.New()
	rout := v1.CreateRouter(cfg, msdb, nil, de)

	srv := server.NewServer(cfg, rout)
	srv.Serve()

}
