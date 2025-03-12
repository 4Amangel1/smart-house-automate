package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/4Amangel1/smart-house-automate/internal/api"
	"github.com/4Amangel1/smart-house-automate/internal/config"
	"github.com/4Amangel1/smart-house-automate/internal/database"
)

func main() {
	logger := log.New(os.Stdout, "smart-house-api: ", log.LstdFlags)
	logger.Println("Starting Smart House API server...")

	// Загрузка конфигурации
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к БД
	repo, err := database.NewRepository(cfg.Database.ConnectionString())
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.Close()
	logger.Println("Connected to database")

	// Создание и запуск API сервера
	apiServer := api.NewServer(repo, logger, cfg.API)

	// Запускаем сервер в отдельной горутине
	go func() {
		logger.Printf("API server starting on port %s", cfg.API.Port)
		if err := apiServer.Run(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start API server: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Printf("Received signal %s, shutting down API server...", sig)

	// Корректное завершение сервера с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(ctx); err != nil {
		logger.Fatalf("Error during server shutdown: %v", err)
	}

	logger.Println("API server shutdown complete")
}
