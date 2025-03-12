package collector

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/4Amangel1/smart-house-automate/internal/database"
	"github.com/4Amangel1/smart-house-automate/internal/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Collector struct {
	sensors  []models.Sensor
	repo     *database.Repository
	logger   *log.Logger
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func New(sensors []models.Sensor, repo *database.Repository, logger *log.Logger) *Collector {
	return &Collector{
		sensors:  sensors,
		repo:     repo,
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

func (c *Collector) Start() {
	c.logger.Printf("Starting collector for %d sensors", len(c.sensors))

	// Запуск HTTP-сервера для метрик Prometheus
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		c.logger.Printf("Starting metrics server on :9090")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			c.logger.Printf("Error starting metrics server: %v", err)
		}
	}()

	// Остальной код запуска коллектора...
	for _, sensor := range c.sensors {
		c.wg.Add(1)
		go c.collectFromSensor(sensor)
	}
}

func (c *Collector) Stop() {
	c.logger.Println("Stopping collector...")
	close(c.stopChan)
	c.wg.Wait()
	c.logger.Println("Collector stopped")
}

func (c *Collector) collectFromSensor(sensor models.Sensor) {
	defer c.wg.Done()

	ticker := time.NewTicker(sensor.Interval())
	defer ticker.Stop()

	c.logger.Printf("Started collecting from sensor %s (type: %s) with interval %v",
		sensor.ID(), sensor.Type(), sensor.Interval())

	for {
		select {
		case <-ticker.C:
			reading, err := sensor.Read()
			if err != nil {
				if err != models.ErrNotReady {
					c.logger.Printf("Error reading from sensor %s: %v", sensor.ID(), err)
				}
				continue
			}

			sensorData := models.SensorData{
				SensorID:   sensor.ID(),
				SensorType: sensor.Type(),
				Timestamp:  reading.Timestamp,
				Value:      models.SensorValue{Data: reading.Value},
				CreatedAt:  time.Now().UTC(),
			}

			if err := c.repo.SaveReading(sensorData); err != nil {
				c.logger.Printf("Error saving reading from sensor %s: %v", sensor.ID(), err)
				readingErrors.Inc() // Увеличиваем счетчик ошибок
				continue
			}

			// Обновляем метрики Prometheus
			updateMetrics(sensorData)

			c.logger.Printf("Collected and saved data from %s sensor %s", sensor.Type(), sensor.ID())

		case <-c.stopChan:
			c.logger.Printf("Stopping collection from sensor %s", sensor.ID())
			return
		}
	}
}
