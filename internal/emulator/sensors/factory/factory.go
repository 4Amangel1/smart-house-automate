package factory

import (
	"fmt"

	"github.com/4Amangel1/smart-house-automate/internal/emulator/sensors/airquality"
	"github.com/4Amangel1/smart-house-automate/internal/emulator/sensors/motion"
	"github.com/4Amangel1/smart-house-automate/internal/emulator/sensors/temperature"
	"github.com/4Amangel1/smart-house-automate/internal/models"
)

type Config struct {
	Temperature []temperature.Config `yaml:"temperature"`
	Motion      []motion.Config      `yaml:"motion"`
	AirQuality  []airquality.Config  `yaml:"air_quality"`
}

type SensorFactory struct {
	sensors map[string]models.Sensor
}

func New(cfg Config) (*SensorFactory, error) {
	f := &SensorFactory{
		sensors: make(map[string]models.Sensor),
	}

	for _, c := range cfg.Temperature {
		s, err := temperature.New(c)
		if err != nil {
			return nil, fmt.Errorf("temperature sensor error: %w", err)
		}
		f.sensors[s.ID()] = s
	}

	for _, c := range cfg.Motion {
		s, err := motion.New(c)
		if err != nil {
			return nil, fmt.Errorf("motion sensor error: %w", err)
		}
		f.sensors[s.ID()] = s
	}

	for _, c := range cfg.AirQuality {
		s, err := airquality.New(c)
		if err != nil {
			return nil, fmt.Errorf("air quality sensor error: %w", err)
		}
		f.sensors[s.ID()] = s
	}

	return f, nil
}

func (f *SensorFactory) GetAllSensors() []models.Sensor {
	sensors := make([]models.Sensor, 0, len(f.sensors))
	for _, s := range f.sensors {
		sensors = append(sensors, s)
	}
	return sensors
}
