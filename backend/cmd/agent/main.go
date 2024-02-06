package main

import (
	"github.com/k6mil6/distributed-calculator/backend/internal/config"
	"github.com/k6mil6/distributed-calculator/backend/pkg/logger"
	"log/slog"
	"time"
)

func main() {
	cfg := config.Get()
	log := logger.SetupLogger(cfg.Env)

	time.Sleep(5 * time.Second)

	log = log.With(slog.String("env", cfg.Env))

	log.Debug("logger debug mode enabled")

}
