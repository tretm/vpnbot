package v1

import (
	"vpngigabot/internal/config"
	"vpngigabot/internal/server/api/v1/routers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func CreateRouter(cfg *config.Config, sm routers.SotrageManager, vpn routers.Vpn, de routers.DataExchenge) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	router.Get("/", routers.Index(cfg, sm))

	router.Post("/", routers.Succes(cfg, sm, de))
	// Создать нового клиента
	// Отобразить форму для заполнени нужных полей
	// router.Get("/{id}", routers.GetConfig(cfg, sm, vpn))
	// router.Get("/ym/{id}", routers.SuccesPay(cfg, sm, de))
	// router.Post("/ex", routers.SendMsgToBot(cfg, sm, de))

	return router
}
