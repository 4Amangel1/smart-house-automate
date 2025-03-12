package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config содержит все настройки приложения
type Config struct {
	Sensors     SensorsConfig `yaml:"sensors"`
	Database    DatabaseConfig
	API         APIConfig
	TelegramBot TelegramBotConfig
}

// DatabaseConfig содержит настройки базы данных
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ConnectionString возвращает строку подключения к БД
func (c DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

// APIConfig содержит настройки API-сервера
type APIConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// TelegramBotConfig содержит настройки Telegram-бота
type TelegramBotConfig struct {
	Token                     string
	AuthorizedUserIDs         []int64
	AlertCheckInterval        time.Duration
	TemperatureAlertThreshold float64
}

// LoadConfig загружает конфигурацию из файла и переменных окружения
func LoadConfig(filePath string) (*Config, error) {
	// Загружаем .env файл, если он существует
	godotenv.Load()

	// Загружаем базовую конфигурацию из файла
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode configuration: %w", err)
	}

	// Загружаем настройки базы данных из переменных окружения
	cfg.Database = loadDatabaseConfig()

	// Загружаем настройки API из переменных окружения
	cfg.API = loadAPIConfig()

	// Загружаем настройки Telegram бота из переменных окружения
	cfg.TelegramBot = loadTelegramBotConfig()

	return &cfg, nil
}

// loadDatabaseConfig загружает настройки БД из переменных окружения
func loadDatabaseConfig() DatabaseConfig {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "smarthouse"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}
}

// loadAPIConfig загружает настройки API из переменных окружения
func loadAPIConfig() APIConfig {
	readTimeout, _ := time.ParseDuration(getEnv("API_READ_TIMEOUT", "10s"))
	writeTimeout, _ := time.ParseDuration(getEnv("API_WRITE_TIMEOUT", "10s"))
	idleTimeout, _ := time.ParseDuration(getEnv("API_IDLE_TIMEOUT", "60s"))

	return APIConfig{
		Port:         getEnv("API_PORT", "8080"),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}

// loadTelegramBotConfig загружает настройки Telegram бота из переменных окружения
func loadTelegramBotConfig() TelegramBotConfig {
	userIDsStr := getEnv("AUTHORIZED_USER_IDS", "")
	var userIDs []int64

	if userIDsStr != "" {
		for _, idStr := range strings.Split(userIDsStr, ",") {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err == nil {
				userIDs = append(userIDs, id)
			}
		}
	}

	checkInterval, _ := time.ParseDuration(getEnv("ALERT_CHECK_INTERVAL", "1m"))
	tempThreshold, _ := strconv.ParseFloat(getEnv("TEMPERATURE_ALERT_THRESHOLD", "30"), 64)

	return TelegramBotConfig{
		Token:                     getEnv("TELEGRAM_BOT_TOKEN", ""),
		AuthorizedUserIDs:         userIDs,
		AlertCheckInterval:        checkInterval,
		TemperatureAlertThreshold: tempThreshold,
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
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
		return fmt.Errorf("sensor ID cannot be empty")
	}
	if c.Min >= c.Max {
		return fmt.Errorf("minimum temperature must be less than maximum")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("polling interval must be positive")
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
		return fmt.Errorf("sensor ID cannot be empty")
	}
	if c.DetectionInterval <= 0 {
		return fmt.Errorf("detection interval must be positive")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("polling interval must be positive")
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
		return fmt.Errorf("sensor ID cannot be empty")
	}
	if c.MinCO2 >= c.MaxCO2 {
		return fmt.Errorf("minimum CO2 level must be less than maximum")
	}
	if c.MinNH3 >= c.MaxNH3 {
		return fmt.Errorf("minimum NH3 level must be less than maximum")
	}
	if c.Interval <= 0 {
		return fmt.Errorf("polling interval must be positive")
	}
	return nil
}