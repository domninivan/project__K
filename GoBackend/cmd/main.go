package main

import (
	"context"
	"log"

	"github.com/domninivan/project__K/GoBackend/internal/api"
	"github.com/domninivan/project__K/GoBackend/internal/repository"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Файл .env не загружен: %v", err)
	}

	ctx := context.Background()

	pool, err := repository.NewPostgres(ctx)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer pool.Close()

	router := api.SetupRouter(pool)

	log.Println("Сервер запущен на http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
