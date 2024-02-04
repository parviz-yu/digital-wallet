package main

import (
	"context"
	"os"

	"github.com/parviz-yu/digital-wallet/api"
	"github.com/parviz-yu/digital-wallet/api/handlers"
	"github.com/parviz-yu/digital-wallet/internal/config"
	"github.com/parviz-yu/digital-wallet/internal/service"
	"github.com/parviz-yu/digital-wallet/internal/storage/postgres"
	"github.com/parviz-yu/digital-wallet/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.NewLogger(cfg.Env)

	strg, err := postgres.NewStorage(context.Background(), cfg)
	if err != nil {
		log.Error("failed to init storage", logger.Error(err))
		os.Exit(1)
	}
	defer strg.CloseDB()

	svc := service.NewService(cfg, log, strg)

	hand := handlers.NewHandler(cfg, log, svc)
	router := api.SetUpRouter(hand, log)

}
