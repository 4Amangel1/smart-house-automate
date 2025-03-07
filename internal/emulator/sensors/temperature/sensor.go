package temperature

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/4Amangel1/smart-house-automate/internal/config"
	"github.com/4Amangel1/smart-house-automate/internal/models"
)

// Sensor представляет датчик температуры
type Sensor struct {
	id       string
	min      float64
	max      float64
	interval time.Duration
}

// NewSensor создает новый датчик температуры
func NewSensor(cfg config.TemperatureConfig) (*Sensor, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("некорректная конфигурация датчика температуры: %w", err)
	}

	return &Sensor{
		id:       cfg.ID,
		min:      cfg.Min,
		max:      cfg.Max,
		interval: cfg.Interval,
	}, nil
}

// ID возвращает идентификатор датчика
func (s *Sensor) ID() string {
	return s.id
}

// Type возвращает тип датчика
func (s *Sensor) Type() string {
	return "temperature"
}

// Interval возвращает интервал опроса датчика
func (s *Sensor) Interval() time.Duration {
	return s.interval
}

// Read считывает показания датчика
func (s *Sensor) Read() (models.Reading, error) {
	// Генерируем случайное значение в диапазоне
	value := s.min + rand.Float64()*(s.max-s.min)
	value = float64(int(value*10)) / 10 // Округляем до 1 десятичного знака

	return models.Reading{
		Value:     value,
		Timestamp: time.Now().UTC(),
	}, nil
}
