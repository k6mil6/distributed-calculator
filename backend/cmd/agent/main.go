package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/k6mil6/distributed-calculator/backend/internal/config"
	"github.com/k6mil6/distributed-calculator/backend/pkg/logger"
	"log/slog"
)

func main() {
	cfg := config.Get()
	log := logger.SetupLogger(cfg.Env)

	log = log.With(slog.String("env", cfg.Env))

	router := chi.NewRouter()

	//router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	log.Debug("logger debug mode enabled")

}
