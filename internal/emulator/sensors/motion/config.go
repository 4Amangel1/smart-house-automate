package motion

import (
	"errors"
	"time"
)

type Config struct {
	ID                string `yaml:"id"`
	Interval          time.Duration
	DetectionInterval time.Duration `yaml:"detection_interval"`
}

func (c Config) Validate() error {
	if c.ID == "" {
		return errors.New("sensor ID is required")
	}
	if c.DetectionInterval < 0 {
		return errors.New("DetectionInterval can't be below 0")
	}
	return nil
}
