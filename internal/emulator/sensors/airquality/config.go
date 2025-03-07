package airquality

import (
	"errors"
	"time"
)

type Config struct {
	ID       string        `yaml:"id"`
	MinCO2   float64       `yaml:"min_co2"`
	MaxCO2   float64       `yaml:"max_co2"`
	MinNH3   float64       `yaml:"min_nh3"`
	MaxNH3   float64       `yaml:"max_nh3"`
	Interval time.Duration `yaml:"interval"`
}

func (c Config) Validate() error {
	if c.ID == "" {
		return errors.New("sensor ID is required")
	}
	if c.MinCO2 >= c.MaxCO2 {
		return errors.New("invalid C02 range ")
	}
	if c.MinNH3 >= c.MaxNH3 {
		return errors.New("invalid NH3 range ")
	}
	return nil
}
