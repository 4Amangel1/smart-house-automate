package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Sensors SensorsConfig `yaml:"sensors"`
}

type SensorsConfig struct {
	Temperature []TemperatureConfig `yaml:"temperature"`
	Motion      []MotionConfig      `yaml:"motion"`
	AirQuality  []AirQualityConfig  `yaml:"air_quality"`
}

type TemperatureConfig struct {
	ID       string        `yaml:"id"`
	Min      float64       `yaml:"min"`
	Max      float64       `yaml:"max"`
	Interval time.Duration `yaml:"interval"`
}

func (c TemperatureConfig) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("ID датчика не может быть пустым")
	}
	if c.Min >= c.Max {
		return fmt.Errorf("минимальная температура должна быть меньше максимальной")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("интервал опроса должен быть положительным")
	}
	return nil
}

type MotionConfig struct {
	ID                string        `yaml:"id"`
	DetectionInterval time.Duration `yaml:"detection_interval"`
	Interval          time.Duration `yaml:"interval"`
}

func (c MotionConfig) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("ID датчика не может быть пустым")
	}
	if c.DetectionInterval <= 0 {
		return fmt.Errorf("интервал обнаружения должен быть положительным")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("интервал опроса должен быть положительным")
	}
	return nil
}

type AirQualityConfig struct {
	ID       string        `yaml:"id"`
	MinCO2   int           `yaml:"min_co2"`
	MaxCO2   int           `yaml:"max_co2"`
	MinNH3   int           `yaml:"min_nh3"`
	MaxNH3   int           `yaml:"max_nh3"`
	Interval time.Duration `yaml:"interval"`
}

func (c AirQualityConfig) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("ID датчика не может быть пустым")
	}
	if c.MinCO2 >= c.MaxCO2 {
		return fmt.Errorf("минимальный уровень CO2 должен быть меньше максимального")
	}
	if c.MinNH3 >= c.MaxNH3 {
		return fmt.Errorf("минимальный уровень NH3 должен быть меньше максимального")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("интервал опроса должен быть положительным")
	}
	return nil
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл конфигурации: %w", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("не удалось декодировать конфигурацию: %w", err)
	}

	return &cfg, nil
}
