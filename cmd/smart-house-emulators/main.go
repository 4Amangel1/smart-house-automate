package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/4Amangel1/smart-house-automate/internal/config"
	"github.com/4Amangel1/smart-house-automate/internal/emulator/sensors/factory"
	"github.com/4Amangel1/smart-house-automate/internal/models"
)

func main() {
	logger := log.New(os.Stdout, "sensor-emulator: ", log.LstdFlags)
	logger.Println("Starting sensor emulator...")

	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Создаем фабрику датчиков
	sensorFactory, err := factory.New(cfg.Sensors)
	if err != nil {
		logger.Fatalf("Failed to create sensor factory: %v", err)
	}

	// Получаем все датчики
	allSensors := sensorFactory.GetAllSensors()

	if len(allSensors) == 0 {
		logger.Fatal("No sensors configured")
	}

	logger.Printf("Initialized %d sensors", len(allSensors))

	done := make(chan struct{})
	var wg sync.WaitGroup

	for _, sensor := range allSensors {
		wg.Add(1)
		go func(s models.Sensor) {
			defer wg.Done()

			ticker := time.NewTicker(s.Interval())
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					reading, err := s.Read()
					if err != nil {
						if err != models.ErrNotReady {
							logger.Printf("Sensor %s error: %v", s.ID(), err)
						}
						continue
					}
					logger.Printf("[%s] %s: %+v", s.Type(), s.ID(), reading.Value)
				case <-done:
					logger.Printf("Stopping sensor %s", s.ID())
					return
				}
			}
		}(sensor)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	<-stopChan
	logger.Println("Received shutdown signal...")

	close(done)

	wg.Wait()
	logger.Println("All sensors stopped. Shutting down.")
}
