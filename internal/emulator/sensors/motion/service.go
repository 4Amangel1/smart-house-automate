package motion

import (
	"emulator/internal/domain/models"
	"math/rand"
	"time"
)

type Service struct {
	id                string
	detectionInterval time.Duration
	lastDetection     time.Time
	interval          time.Duration
	rnd               *rand.Rand
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func New(cfg Config) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		id:                cfg.ID,
		detectionInterval: cfg.DetectionInterval,
		interval:          cfg.Interval,
		rnd:               rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func (s *Service) Interval() time.Duration {
	return s.interval
}

func (s *Service) ID() string   { return s.id }
func (s *Service) Type() string { return "motion" }

func (s *Service) Read() (models.Reading, error) {
	now := time.Now().UTC()

	if now.Sub(s.lastDetection) < s.detectionInterval {
		return models.Reading{}, models.ErrNotReady
	}

	s.lastDetection = now
	detected := s.rnd.Intn(2) == 1

	var valueDisplay string
	if detected {
		valueDisplay = "motion detected"
	} else {
		valueDisplay = "motion not detected"
	}

	return models.Reading{
		Value:     valueDisplay,
		Timestamp: now,
	}, nil
}
