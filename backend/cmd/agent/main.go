package main

import (
	"github.com/k6mil6/distributed-calculator/backend/internal/agent/worker"
	"github.com/k6mil6/distributed-calculator/backend/internal/config"
	"github.com/k6mil6/distributed-calculator/backend/pkg/logger"
	"log/slog"
)

func main() {
	cfg := config.Get()
	log := logger.SetupLogger(cfg.Env)

	log = log.With(slog.String("env", cfg.Env))

	for i := 0; i < cfg.AgentsNumber; i++ {
		go func(i int) {
			w := worker.New(int64(i))
			w.Start(cfg.OrchestratorURL, log, cfg.WorkerTimeout, cfg.HeartbeatTimeout)
		}(i)
	}

	log.Debug("logger debug mode enabled")

}
