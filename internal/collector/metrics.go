package collector

import (
	"strings"

	"github.com/4Amangel1/smart-house-automate/internal/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Метрики температуры
	temperatureGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "temperature_celsius",
			Help: "Температура в градусах Цельсия",
		},
		[]string{"sensor_id", "location"},
	)

	// Метрики движения
	motionGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "motion_detected",
			Help: "Обнаружено ли движение (1 - да, 0 - нет)",
		},
		[]string{"sensor_id", "location"},
	)

	// Метрики качества воздуха
	co2Gauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "air_co2_ppm",
			Help: "Уровень CO2 в частях на миллион",
		},
		[]string{"sensor_id", "location"},
	)

	nh3Gauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "air_nh3_ppm",
			Help: "Уровень NH3 в частях на миллион",
		},
		[]string{"sensor_id", "location"},
	)

	// Общие метрики
	readingsCollected = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "collector_readings_total",
			Help: "Общее количество собранных показаний",
		},
	)

	readingErrors = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "collector_reading_errors_total",
			Help: "Общее количество ошибок при сборе показаний",
		},
	)
)

// Обновление метрик при сборе данных
func updateMetrics(data models.SensorData) {
	readingsCollected.Inc()

	// Получите местоположение из ID (например, temp_living_room => living_room)
	parts := strings.Split(data.SensorID, "_")
	location := "unknown"
	if len(parts) > 1 {
		location = strings.Join(parts[1:], "_")
	}

	switch data.SensorType {
	case "temperature":
		if value, ok := data.Value.Data.(float64); ok {
			temperatureGauge.WithLabelValues(data.SensorID, location).Set(value)
		}

	case "motion":
		if value, ok := data.Value.Data.(bool); ok {
			if value {
				motionGauge.WithLabelValues(data.SensorID, location).Set(1)
			} else {
				motionGauge.WithLabelValues(data.SensorID, location).Set(0)
			}
		}

	case "air_quality":
		if valueMap, ok := data.Value.Data.(map[string]interface{}); ok {
			if co2, co2ok := valueMap["co2"].(float64); co2ok {
				co2Gauge.WithLabelValues(data.SensorID, location).Set(co2)
			}
			if nh3, nh3ok := valueMap["nh3"].(float64); nh3ok {
				nh3Gauge.WithLabelValues(data.SensorID, location).Set(nh3)
			}
		}
	}
}
