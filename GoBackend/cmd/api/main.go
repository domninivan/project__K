package main

import (
	"context"
	"log"

	"github.com/domninivan/project__K/GoBackend/internal/config"
	"github.com/domninivan/project__K/GoBackend/internal/handler"
	"github.com/domninivan/project__K/GoBackend/internal/repository"
	router "github.com/domninivan/project__K/GoBackend/internal/router"
	"github.com/domninivan/project__K/GoBackend/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	pool, err := repository.NewPostgresPool(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	healthRepo := repository.NewHealth(pool)
	workCenterRepo := repository.NewWorkCenter(pool)

	healthSvc := service.NewHealth(healthRepo)
	workCenterSvc := service.NewWorkCenter(workCenterRepo)

	handlers := router.Handlers{
		Health:     handler.NewHealth(healthSvc),
		WorkCenter: handler.NewWorkCenter(workCenterSvc),
	}

	engine := router.New(handlers)

	log.Printf("Server started at http://localhost:%s", cfg.ServerPort)
	if err := engine.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
