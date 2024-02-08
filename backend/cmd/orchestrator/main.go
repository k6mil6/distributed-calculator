package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/backend/internal/config"
	"github.com/k6mil6/distributed-calculator/backend/pkg/logger"
	"log/slog"
	"time"
)

func main() {
	cfg := config.Get()
	log := logger.SetupLogger(cfg.Env).With(slog.String("env", cfg.Env))

	db, err := sqlx.Connect("postgres", cfg.DatabaseDSN)
	if err != nil {
		log.Error("failed to connect to database", err)
		return
	}
	defer db.Close()

	//expressionStorage := storage.NewExpressionStorage(db)
	//
	//expressionStorage.NonTakenExpressions(context.Background())

	time.Sleep(5 * time.Second)

	log.Debug("logger debug mode enabled")
}
