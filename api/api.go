package api

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/parviz-yu/digital-wallet/api/handlers"
	"github.com/parviz-yu/digital-wallet/pkg/logger"
)

func SetUpRouter(h *handlers.Handler, log logger.LoggerI) *chi.Mux {
	router := chi.NewRouter()

	router.Use(handlers.NewMWLogger(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.RequestID)
	router.Use(handlers.AuthMiddlewareUserID)

	router.Head("/api/v1/wallets", h.DoesWalletExists)
	router.Post("/api/v1/wallets", h.PutFunds())
	router.Get("/api/v1/wallets/stats", h.GetStats)
	router.Get("/api/v1/wallets/balance", h.GetBalance)

	return router
}
