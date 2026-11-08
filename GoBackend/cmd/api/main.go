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

	wcRepo := repository.NewWorkCenterRepository(pool)
	wcSvc := service.NewWorkCenterService(wcRepo)
	wcH := handler.NewWorkCenterHandler(wcSvc)

	handlers := router.Handlers{
		WorkCenter: wcH,
	}

	engine := router.SetupRouter(pool, handlers)

	log.Printf("Server started at http://localhost:%s", cfg.ServerPort)
	if err := engine.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
