package logger

import (
	"compile-server/internal/logger/slogpretty"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

const (
	Err = "error"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slogpretty.SetupSlogPretty()
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
