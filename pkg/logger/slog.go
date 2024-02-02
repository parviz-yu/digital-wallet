package logger

import (
	"log"
	"os"

	"golang.org/x/exp/slog"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

type loggerWrapper struct {
	slog *slog.Logger
}

func newSlogLogger(level string) *slog.Logger {
	var loggger *slog.Logger

	switch level {
	case envLocal:
		loggger = slog.New(
			slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		loggger = slog.New(
			slog.NewJSONHandler(os.Stdin, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log.Fatalf("invalid environment: %s", level)
	}

	return loggger
}
