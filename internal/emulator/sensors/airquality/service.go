package airquality

import (
	"math/rand"
	"time"

	"github.com/4Amangel1/smart-house-automate/internal/models"
)

type Service struct {
	id       string
	minCO2   float64
	maxCO2   float64
	minNH3   float64
	maxNH3   float64
	interval time.Duration
}

func New(cfg Config) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		id:       cfg.ID,
		minCO2:   cfg.MinCO2,
		maxCO2:   cfg.MaxCO2,
		minNH3:   cfg.MinNH3,
		maxNH3:   cfg.MaxNH3,
		interval: cfg.Interval,
	}, nil

}

func (s *Service) Interval() time.Duration {
	return s.interval
}

func (s *Service) ID() string   { return s.id }
func (s *Service) Type() string { return "air_quality" }

func (s *Service) Read() (models.Reading, error) {
	co2 := s.minCO2 + rand.Float64()*(s.maxCO2-s.minCO2)
	nh3 := s.minNH3 + rand.Float64()*(s.maxNH3-s.minNH3)

	return models.Reading{
		Value: struct {
			CO2     float64 `json:"co2"`
			Ammonia float64 `json:"ammonia"`
		}{
			CO2:     co2,
			Ammonia: nh3,
		},
		Timestamp: time.Now().UTC(),
	}, nil
}
