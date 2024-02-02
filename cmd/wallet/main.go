package main

import (
	"github.com/parviz-yu/digital-wallet/internal/config"
	"github.com/parviz-yu/digital-wallet/pkg/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.NewLogger(cfg.Env)

	// init storage

}
