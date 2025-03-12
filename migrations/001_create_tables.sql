-- Таблица для хранения показаний датчиков
CREATE TABLE IF NOT EXISTS sensor_readings (
    id SERIAL PRIMARY KEY,
    sensor_id VARCHAR(100) NOT NULL,
    sensor_type VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    value JSONB NOT NULL,
    unit VARCHAR(20),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_sensor_readings_sensor_id ON sensor_readings(sensor_id);
CREATE INDEX IF NOT EXISTS idx_sensor_readings_timestamp ON sensor_readings(timestamp);
CREATE INDEX IF NOT EXISTS idx_sensor_readings_sensor_type ON sensor_readings(sensor_type);
CREATE INDEX IF NOT EXISTS idx_sensor_readings_created_at ON sensor_readings(created_at);

-- Таблица для хранения информации о датчиках
CREATE TABLE IF NOT EXISTS sensors (
    id VARCHAR(100) PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    name VARCHAR(100),
    location VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
); 