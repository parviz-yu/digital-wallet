package main

import (
	"context"
	"os"

	"github.com/parviz-yu/digital-wallet/internal/config"
	"github.com/parviz-yu/digital-wallet/internal/storage/postgres"
	"github.com/parviz-yu/digital-wallet/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.NewLogger(cfg.Env)

	// init storage
	strg, err := postgres.NewStorage(context.Background(), cfg)
	if err != nil {
		log.Error("failed to init storage", logger.Error(err))
		os.Exit(1)
	}

	// init services

	//
}
