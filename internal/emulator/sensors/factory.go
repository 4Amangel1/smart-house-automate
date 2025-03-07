package sensors

import (
	"github.com/lostly/smart-house-automate/internal/config"
	"github.com/lostly/smart-house-automate/internal/domain/models"
	"github.com/lostly/smart-house-automate/internal/emulator/sensors/airquality"
	"github.com/lostly/smart-house-automate/internal/emulator/sensors/motion"
	"github.com/lostly/smart-house-automate/internal/emulator/sensors/temperature"
)

// NewTemperatureSensor создает новый датчик температуры
func NewTemperatureSensor(cfg config.TemperatureConfig) (models.Sensor, error) {
	return temperature.NewSensor(cfg)
}

// NewMotionSensor создает новый датчик движения
func NewMotionSensor(cfg config.MotionConfig) (models.Sensor, error) {
	return motion.NewSensor(cfg)
}

// NewAirQualitySensor создает новый датчик качества воздуха
func NewAirQualitySensor(cfg config.AirQualityConfig) (models.Sensor, error) {
	return airquality.NewSensor(cfg)
}
