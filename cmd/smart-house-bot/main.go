package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/4Amangel1/smart-house-automate/internal/bot"
	"github.com/4Amangel1/smart-house-automate/internal/config"
	"github.com/4Amangel1/smart-house-automate/internal/database"
)

func main() {
	logger := log.New(os.Stdout, "smart-house-bot: ", log.LstdFlags)
	logger.Println("Starting Smart House Telegram Bot...")

	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	if cfg.TelegramBot.Token == "" {
		logger.Fatal("Telegram bot token is not set")
	}

	repo, err := database.NewRepository(cfg.Database.ConnectionString())
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.Close()
	logger.Println("Connected to database")

	tgBot, err := bot.NewBot(cfg.TelegramBot, repo, logger)
	if err != nil {
		logger.Fatalf("Failed to create Telegram bot: %v", err)
	}

	go func() {
		logger.Println("Starting Telegram bot...")
		if err := tgBot.Start(); err != nil {
			logger.Fatalf("Telegram bot error: %v", err)
		}
	}()

	go func() {
		logger.Println("Starting notification service...")
		tgBot.StartNotificationService()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Printf("Received signal %s, shutting down Telegram bot...", sig)

	tgBot.Stop()

	time.Sleep(1 * time.Second)

	logger.Println("Telegram bot shutdown complete")
}
