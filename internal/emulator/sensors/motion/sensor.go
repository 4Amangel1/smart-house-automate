package motion

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lostly/smart-house-automate/internal/config"
	"github.com/lostly/smart-house-automate/internal/domain/models"
)

// Sensor представляет датчик движения
type Sensor struct {
	id                string
	detectionInterval time.Duration
	interval          time.Duration
	lastDetected      time.Time
}

// NewSensor создает новый датчик движения
func NewSensor(cfg config.MotionConfig) (*Sensor, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("некорректная конфигурация датчика движения: %w", err)
	}

	return &Sensor{
		id:                cfg.ID,
		detectionInterval: cfg.DetectionInterval,
		interval:          cfg.Interval,
		lastDetected:      time.Now().Add(-24 * time.Hour), // Инициализируем с прошлым днем
	}, nil
}

// ID возвращает идентификатор датчика
func (s *Sensor) ID() string {
	return s.id
}

// Type возвращает тип датчика
func (s *Sensor) Type() string {
	return "motion"
}

// Interval возвращает интервал опроса датчика
func (s *Sensor) Interval() time.Duration {
	return s.interval
}

// Read считывает показания датчика
func (s *Sensor) Read() (models.Reading, error) {
	now := time.Now()

	// 30% шанс обнаружения движения
	detected := rand.Float64() < 0.3

	// Если движение обнаружено, обновляем время последнего обнаружения
	if detected {
		s.lastDetected = now
	}

	// Движение считается активным, если оно было обнаружено в течение detectionInterval
	isActive := now.Sub(s.lastDetected) <= s.detectionInterval

	return models.Reading{
		Value:     isActive,
		Timestamp: now.UTC(),
	}, nil
}
