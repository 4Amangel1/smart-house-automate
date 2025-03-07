package airquality

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/4Amangel1/smart-house-automate/internal/config"
	"github.com/4Amangel1/smart-house-automate/internal/models"
)

// Sensor представляет датчик качества воздуха
type Sensor struct {
	id       string
	minCO2   float64
	maxCO2   float64
	minNH3   float64
	maxNH3   float64
	interval time.Duration
}

// NewSensor создает новый датчик качества воздуха
func NewSensor(cfg config.AirQualityConfig) (*Sensor, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("некорректная конфигурация датчика качества воздуха: %w", err)
	}

	return &Sensor{
		id:       cfg.ID,
		minCO2:   float64(cfg.MinCO2),
		maxCO2:   float64(cfg.MaxCO2),
		minNH3:   float64(cfg.MinNH3),
		maxNH3:   float64(cfg.MaxNH3),
		interval: cfg.Interval,
	}, nil
}

// ID возвращает идентификатор датчика
func (s *Sensor) ID() string {
	return s.id
}

// Type возвращает тип датчика
func (s *Sensor) Type() string {
	return "air_quality"
}

// Interval возвращает интервал опроса датчика
func (s *Sensor) Interval() time.Duration {
	return s.interval
}

// Read считывает показания датчика
func (s *Sensor) Read() (models.Reading, error) {
	// Генерируем случайные значения в диапазоне
	co2 := s.minCO2 + rand.Float64()*(s.maxCO2-s.minCO2)
	nh3 := s.minNH3 + rand.Float64()*(s.maxNH3-s.minNH3)

	// Округляем до целых значений
	co2 = float64(int(co2))
	nh3 = float64(int(nh3))

	data := map[string]interface{}{
		"co2": co2,
		"nh3": nh3,
	}

	return models.Reading{
		Value:     data,
		Timestamp: time.Now().UTC(),
	}, nil
}
