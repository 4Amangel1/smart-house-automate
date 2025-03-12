package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/4Amangel1/smart-house-automate/internal/collector"
	"github.com/4Amangel1/smart-house-automate/internal/config"
	"github.com/4Amangel1/smart-house-automate/internal/database"
	"github.com/4Amangel1/smart-house-automate/internal/emulator/sensors/factory"
)

func main() {
	logger := log.New(os.Stdout, "smart-house-collector: ", log.LstdFlags)
	logger.Println("Starting Smart House Collector service...")

	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	repo, err := database.NewRepository(cfg.Database.ConnectionString())
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.Close()
	logger.Println("Connected to database")

	sensorFactory, err := factory.New(cfg.Sensors)
	if err != nil {
		logger.Fatalf("Failed to create sensor factory: %v", err)
	}

	allSensors := sensorFactory.GetAllSensors()
	if len(allSensors) == 0 {
		logger.Fatal("No sensors configured")
	}
	logger.Printf("Initialized %d sensors", len(allSensors))

	dataCollector := collector.New(allSensors, repo, logger)
	dataCollector.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	logger.Printf("Received signal %s, shutting down collector...", sig)

	dataCollector.Stop()

	logger.Println("Collector shutdown complete")
}
