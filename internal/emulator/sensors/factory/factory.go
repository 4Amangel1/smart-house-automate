package factory

import (
	"fmt"

	"github.com/4Amangel1/smart-house-automate/internal/config"
	"github.com/4Amangel1/smart-house-automate/internal/emulator/sensors/airquality"
	"github.com/4Amangel1/smart-house-automate/internal/emulator/sensors/motion"
	"github.com/4Amangel1/smart-house-automate/internal/emulator/sensors/temperature"
	"github.com/4Amangel1/smart-house-automate/internal/models"
)

type SensorFactory struct {
	sensors map[string]models.Sensor
}

func New(cfg config.SensorsConfig) (*SensorFactory, error) {
	f := &SensorFactory{
		sensors: make(map[string]models.Sensor),
	}

	// Температурные датчики
	for _, c := range cfg.Temperature {
		// Преобразуем общую конфигурацию в конфигурацию для датчика
		sensorCfg := temperature.Config{
			ID:       c.ID,
			Min:      c.Min,
			Max:      c.Max,
			Interval: c.Interval,
		}

		s, err := temperature.New(sensorCfg)
		if err != nil {
			return nil, fmt.Errorf("temperature sensor error: %w", err)
		}
		f.sensors[s.ID()] = s
	}

	// Датчики движения
	for _, c := range cfg.Motion {
		sensorCfg := motion.Config{
			ID:                c.ID,
			DetectionInterval: c.DetectionInterval,
			Interval:          c.Interval,
		}

		s, err := motion.New(sensorCfg)
		if err != nil {
			return nil, fmt.Errorf("motion sensor error: %w", err)
		}
		f.sensors[s.ID()] = s
	}

	// Датчики качества воздуха
	for _, c := range cfg.AirQuality {
		sensorCfg := airquality.Config{
			ID:       c.ID,
			MinCO2:   float64(c.MinCO2),
			MaxCO2:   float64(c.MaxCO2),
			MinNH3:   float64(c.MinNH3),
			MaxNH3:   float64(c.MaxNH3),
			Interval: c.Interval,
		}

		s, err := airquality.New(sensorCfg)
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
