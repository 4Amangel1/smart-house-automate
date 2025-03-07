package temperature

import (
	"math/rand"
	"time"

	"github.com/4Amangel1/smart-house-automate/internal/models"
)

type Service struct {
	id       string
	min      float64
	max      float64
	interval time.Duration
}

func New(cfg Config) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		id:       cfg.ID,
		min:      cfg.Min,
		max:      cfg.Max,
		interval: cfg.Interval,
	}, nil
}

func (s *Service) Interval() time.Duration {
	return s.interval
}

func (s *Service) ID() string   { return s.id }
func (s *Service) Type() string { return "temperature" }

func (s *Service) Read() (models.Reading, error) {
	value := s.min + rand.Float64()*(s.max-s.min)
	return models.Reading{
		Value:     value,
		Timestamp: time.Now().UTC(),
	}, nil
}
