package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/domninivan/project__K/GoBackend/internal/config"
	"github.com/domninivan/project__K/GoBackend/internal/handler"
	"github.com/domninivan/project__K/GoBackend/internal/logger"
	"github.com/domninivan/project__K/GoBackend/internal/repository"
	"github.com/domninivan/project__K/GoBackend/internal/router"
	"github.com/domninivan/project__K/GoBackend/internal/service"
)

func main() {
	if err := run(); err != nil {
		slog.Default().Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	l, err := logger.New(cfg.LogLevel, cfg.LogFormat, os.Stdout)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	slog.SetDefault(l)

	mainLog := l.With("component", "main")
	dbLog := l.With("component", "db")
	httpLog := l.With("component", "http")

	mainLog.Info("starting")

	ctx := context.Background()

	pool, err := repository.NewPostgresPool(ctx, cfg)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer pool.Close()
	dbLog.Info("connected to database")

	healthRepo := repository.NewHealth(pool)
	workCenterRepo := repository.NewWorkCenter(pool)

	healthSvc := service.NewHealth(healthRepo)
	workCenterSvc := service.NewWorkCenter(workCenterRepo)

	handlers := router.Handlers{
		Health:     handler.NewHealth(healthSvc),
		WorkCenter: handler.NewWorkCenter(workCenterSvc),
	}

	engine := router.New(handlers)

	httpLog.Info("listening", "addr", ":"+cfg.ServerPort)
	if err := engine.Run(":" + cfg.ServerPort); err != nil {
		return fmt.Errorf("start server: %w", err)
	}
	return nil
}
