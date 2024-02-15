package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/distributed-calculator/backend/internal/config"
	"github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/fetcher"
	"github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/http_server/handlers/agents/free_expressions"
	"github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/http_server/handlers/agents/result"
	"github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/http_server/handlers/expression/all_expressions"
	"github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/http_server/handlers/expression/calculate"
	mwlogger "github.com/k6mil6/distributed-calculator/backend/internal/orchestrator/http_server/middleware/logger"
	"github.com/k6mil6/distributed-calculator/backend/internal/storage/migrations"
	"github.com/k6mil6/distributed-calculator/backend/internal/storage/postgres"
	"github.com/k6mil6/distributed-calculator/backend/pkg/logger"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Get()
	log := logger.SetupLogger(cfg.Env).With(slog.String("env", cfg.Env))
	log.Debug("logger debug mode enabled")

	db, err := sqlx.Connect("postgres", cfg.DatabaseDSN)
	if err != nil {
		log.Error("failed to connect to database", err)
		return
	}
	defer db.Close()

	if err := migrations.Start(cfg.DatabaseDSN); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Error("failed to start migrations", err)
			return
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	expressionStorage := postgres.NewExpressionStorage(db)
	subExpressionStorage := postgres.NewSubExpressionStorage(db)

	router := chi.NewRouter()

	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/calculate", calculate.New(log, expressionStorage, ctx))
	router.Post("/free_expressions", free_expressions.New(log, subExpressionStorage, ctx))
	router.Post("/result", result.New(log, subExpressionStorage, ctx))
	router.Get("/all_expressions", all_expressions.New(log, expressionStorage, subExpressionStorage, ctx))

	fetcher := fetcher.New(expressionStorage, subExpressionStorage, cfg.FetcherInterval, log)

	srv := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	go fetcher.Start(ctx)

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", err)
	}

	log.Info("server stopped")

}
