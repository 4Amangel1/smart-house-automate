package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/4Amangel1/smart-house-automate/internal/models"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(connStr string) (*Repository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

// SaveReading сохраняет показания датчика в БД
func (r *Repository) SaveReading(data models.SensorData) error {
	query := `
		INSERT INTO sensor_readings (sensor_id, sensor_type, timestamp, value, unit, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	valueJSON, err := data.Value.Value()
	if err != nil {
		return fmt.Errorf("failed to convert value to JSON: %w", err)
	}

	metadataJSON, err := data.Metadata.Value()
	if err != nil {
		return fmt.Errorf("failed to convert metadata to JSON: %w", err)
	}

	var id int64
	err = r.db.QueryRow(
		query,
		data.SensorID,
		data.SensorType,
		data.Timestamp,
		valueJSON,
		data.Unit,
		metadataJSON,
		time.Now().UTC(),
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("failed to save reading: %w", err)
	}

	return nil
}

// GetLatestReadings возвращает последние показания со всех датчиков
func (r *Repository) GetLatestReadings() ([]models.SensorData, error) {
	query := `
		SELECT DISTINCT ON (sensor_id)
			id, sensor_id, sensor_type, timestamp, value, unit, metadata, created_at
		FROM sensor_readings
		ORDER BY sensor_id, timestamp DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest readings: %w", err)
	}
	defer rows.Close()

	return r.scanReadings(rows)
}

// GetLatestReadingsByType возвращает последние показания с датчиков определенного типа
func (r *Repository) GetLatestReadingsByType(sensorType string) ([]models.SensorData, error) {
	query := `
		SELECT DISTINCT ON (sensor_id)
			id, sensor_id, sensor_type, timestamp, value, unit, metadata, created_at
		FROM sensor_readings
		WHERE sensor_type = $1
		ORDER BY sensor_id, timestamp DESC
	`

	rows, err := r.db.Query(query, sensorType)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest readings by type: %w", err)
	}
	defer rows.Close()

	return r.scanReadings(rows)
}

// GetReadingHistory возвращает историю показаний конкретного датчика
func (r *Repository) GetReadingHistory(sensorID string, limit int) ([]models.SensorData, error) {
	query := `
		SELECT id, sensor_id, sensor_type, timestamp, value, unit, metadata, created_at
		FROM sensor_readings
		WHERE sensor_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, sensorID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query reading history: %w", err)
	}
	defer rows.Close()

	return r.scanReadings(rows)
}

// scanReadings сканирует результаты запроса в структуры моделей
func (r *Repository) scanReadings(rows *sql.Rows) ([]models.SensorData, error) {
	var readings []models.SensorData

	for rows.Next() {
		var reading models.SensorData
		var valueBytes, metadataBytes []byte
		var unit sql.NullString

		if err := rows.Scan(
			&reading.ID,
			&reading.SensorID,
			&reading.SensorType,
			&reading.Timestamp,
			&valueBytes,
			&unit,
			&metadataBytes,
			&reading.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan reading: %w", err)
		}

		// Установка значения
		if err := reading.Value.Scan(valueBytes); err != nil {
			return nil, fmt.Errorf("failed to scan value: %w", err)
		}

		// Установка метаданных
		if err := reading.Metadata.Scan(metadataBytes); err != nil {
			return nil, fmt.Errorf("failed to scan metadata: %w", err)
		}

		// Установка единицы измерения
		if unit.Valid {
			reading.Unit = unit.String
		}

		readings = append(readings, reading)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return readings, nil
}
