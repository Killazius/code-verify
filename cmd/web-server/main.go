package main

import (
	"compile-server/internal/config"
	"compile-server/internal/handlers/ws"
	"compile-server/internal/logger"
	"compile-server/internal/middleware/customLogger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func main() {

	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(customLogger.New(log))
	router.Use(middleware.Recoverer)

	router.HandleFunc("/ws", ws.New(log, cfg.Env))
	server := http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("starting server", "address", cfg.Address, "env", cfg.Env)
	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server", slog.String("error", err.Error()))
	}
}
