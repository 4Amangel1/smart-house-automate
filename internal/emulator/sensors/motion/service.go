package motion

import (
	"math/rand"
	"time"

	"github.com/4Amangel1/smart-house-automate/internal/models"
)

type Service struct {
	id                string
	detectionInterval time.Duration
	lastDetection     time.Time
	interval          time.Duration
}

func New(cfg Config) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		id:                cfg.ID,
		detectionInterval: cfg.DetectionInterval,
		interval:          cfg.Interval,
		lastDetection:     time.Now().Add(-24 * time.Hour), // Инициализация в прошлом
	}, nil
}

func (s *Service) Interval() time.Duration {
	return s.interval
}

func (s *Service) ID() string   { return s.id }
func (s *Service) Type() string { return "motion" }

func (s *Service) Read() (models.Reading, error) {
	now := time.Now().UTC()

	// 30% шанс обнаружения движения
	detected := rand.Float64() < 0.3

	// Если движение обнаружено, обновляем время последнего обнаружения
	if detected {
		s.lastDetection = now
	}

	// Движение считается активным, если оно было обнаружено в течение detectionInterval
	isActive := now.Sub(s.lastDetection) <= s.detectionInterval

	return models.Reading{
		Value:     isActive,
		Timestamp: now,
	}, nil
}
