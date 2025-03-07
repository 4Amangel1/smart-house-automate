package temperature

import (
	"errors"
	"time"
)

type Config struct {
	ID       string        `yaml:"id"`
	Min      float64       `yaml:"min"`
	Max      float64       `yaml:"max"`
	Interval time.Duration `yaml:"interval"`
}

func (c Config) Validate() error {
	if c.ID == "" {
		return errors.New("sensor ID is required")
	}
	if c.Min >= c.Max {
		return errors.New("invalid temperature range")
	}
	return nil
}
