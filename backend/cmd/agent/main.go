package main

import (
	"github.com/k6mil6/distributed-calculator/backend/internal/agent/worker"
	"github.com/k6mil6/distributed-calculator/backend/internal/config"
	"github.com/k6mil6/distributed-calculator/backend/pkg/logger"
	"log/slog"
	"sync"
)

func main() {
	cfg := config.Get()
	log := logger.SetupLogger(cfg.Env)

	log = log.With(slog.String("env", cfg.Env))
	log.Debug("logger debug mode enabled")

	var wg sync.WaitGroup

	for i := 0; i < cfg.GoroutineNumber; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			w := worker.New()
			w.Start(cfg.OrchestratorURL, log, cfg.WorkerTimeout, cfg.HeartbeatTimeout)
		}(i)
	}

	wg.Wait()
}
