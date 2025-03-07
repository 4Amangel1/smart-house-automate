package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Interfaces

// Sensor представляет интерфейс для датчика
type Sensor interface {
	ID() string
	Type() string
	Read() (Reading, error)
	Interval() time.Duration
}

// Reading представляет данные, полученные с датчика
type Reading struct {
	Value     interface{} `json:"value"`
	Timestamp time.Time   `json:"timestamp"`
}

// SensorConfig представляет интерфейс для конфигурации датчика
type SensorConfig interface {
	Validate() error
}

// SensorData представляет данные, полученные от датчика для сохранения в БД
type SensorData struct {
	ID         int64       `json:"id" db:"id"`                       // Уникальный идентификатор записи
	SensorID   string      `json:"sensorId" db:"sensor_id"`          // Идентификатор датчика
	SensorType string      `json:"sensorType" db:"sensor_type"`      // Тип датчика
	Timestamp  time.Time   `json:"timestamp" db:"timestamp"`         // Время получения данных
	Value      SensorValue `json:"value" db:"value"`                 // Значение данных
	Unit       string      `json:"unit,omitempty" db:"unit"`         // Единица измерения
	Metadata   Metadata    `json:"metadata,omitempty" db:"metadata"` // Дополнительные метаданные
	CreatedAt  time.Time   `json:"createdAt" db:"created_at"`        // Время создания записи
}

// SensorValue - JSON-тип для значения датчика
type SensorValue struct {
	Data interface{} `json:"data"`
}

// Value - реализация интерфейса driver.Valuer для записи в БД
func (sv SensorValue) Value() (driver.Value, error) {
	return json.Marshal(sv.Data)
}

// Scan - реализация интерфейса sql.Scanner для чтения из БД
func (sv *SensorValue) Scan(value interface{}) error {
	if value == nil {
		sv.Data = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("тип не []byte")
	}

	var data interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	sv.Data = data
	return nil
}

// Metadata - JSON-тип для метаданных
type Metadata struct {
	Data map[string]interface{} `json:"data"`
}

// Value - реализация интерфейса driver.Valuer для записи в БД
func (m Metadata) Value() (driver.Value, error) {
	if m.Data == nil {
		return nil, nil
	}
	return json.Marshal(m.Data)
}

// Scan - реализация интерфейса sql.Scanner для чтения из БД
func (m *Metadata) Scan(value interface{}) error {
	if value == nil {
		m.Data = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("тип не []byte")
	}

	if len(bytes) == 0 {
		m.Data = nil
		return nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	m.Data = data
	return nil
}

// Errors
var (
	ErrNotReady = errors.New("sensor is not ready to read")
)
