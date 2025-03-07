package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/lostly/smart-house-automate/internal/api"
	"github.com/lostly/smart-house-automate/internal/config"
	"github.com/lostly/smart-house-automate/internal/domain/models"
	"github.com/lostly/smart-house-automate/internal/emulator/sensors"
)

func main() {
	logger := log.New(os.Stdout, "sensor-emulator: ", log.LstdFlags)
	logger.Println("Starting sensor emulator...")

	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Создаем и инициализируем датчики
	allSensors, err := initializeSensors(cfg.Sensors)
	if err != nil {
		logger.Fatalf("Failed to initialize sensors: %v", err)
	}

	// Инициализация API клиента
	apiClient := api.NewClient(cfg.API)
	apiService := api.NewService(apiClient, time.Duration(cfg.API.SendInterval)*time.Second)

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
						logger.Printf("Sensor %s error: %v", s.ID(), err)
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

	// Регистрируем все датчики для отправки данных
	for _, sensor := range allSensors {
		apiService.RegisterSensor(sensor, sensor.Type())
	}

	// Запускаем сервис отправки данных
	apiService.Start()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	<-stopChan
	logger.Println("Received shutdown signal...")

	close(done)

	wg.Wait()
	logger.Println("All sensors stopped. Shutting down.")

	// Останавливаем сервис отправки данных
	apiService.Stop()
}

// initializeSensors инициализирует все датчики на основе конфигурации
func initializeSensors(cfg config.SensorsConfig) ([]models.Sensor, error) {
	var allSensors []models.Sensor

	// Инициализация температурных датчиков
	for _, tCfg := range cfg.Temperature {
		sensor, err := sensors.NewTemperatureSensor(tCfg)
		if err != nil {
			return nil, err
		}
		allSensors = append(allSensors, sensor)
	}

	// Инициализация датчиков движения
	for _, mCfg := range cfg.Motion {
		sensor, err := sensors.NewMotionSensor(mCfg)
		if err != nil {
			return nil, err
		}
		allSensors = append(allSensors, sensor)
	}

	// Инициализация датчиков качества воздуха
	for _, aqCfg := range cfg.AirQuality {
		sensor, err := sensors.NewAirQualitySensor(aqCfg)
		if err != nil {
			return nil, err
		}
		allSensors = append(allSensors, sensor)
	}

	return allSensors, nil
}
