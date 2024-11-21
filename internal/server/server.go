package server

import (
	"fmt"
	"net/http"
	"vpngigabot/internal/config"
)

type Server struct {
	cfg *config.Config
	srv http.Server
}

func NewServer(cfg *config.Config, router http.Handler) *Server {
	srv := Server{
		cfg: cfg,
		srv: http.Server{
			Addr:         ":" + cfg.ServerHttpConfig.Port,
			Handler:      router,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
	return &srv
}

// Serve http request
func (s *Server) Serve() {
	fmt.Println("http server serves on port " + s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Println("server stopped with error: ", err.Error())

	}
}
