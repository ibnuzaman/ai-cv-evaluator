package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Underscore for side-effect import
)

type DBConfig struct {
	Host               string
	Port               int
	User               string
	Password           string
	Name               string
	SSLMode            string
	MaxOpen            int
	MaxIdle            int
	ConnectionLifetime time.Duration
	ConnectionIdle     time.Duration
}

type Config struct {
	AppPort      string
	DB           *DBConfig
	DatabaseURL  string
	GeminiAPIKey string
	ChromaDBURL  string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	// Parse database configuration
	dbPort, err := strconv.Atoi(getEnvOrDefault("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	maxOpen, err := strconv.Atoi(getEnvOrDefault("DB_MAX_OPEN", "25"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_OPEN: %w", err)
	}

	maxIdle, err := strconv.Atoi(getEnvOrDefault("DB_MAX_IDLE", "25"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_IDLE: %w", err)
	}

	connLifetime, err := time.ParseDuration(getEnvOrDefault("DB_CONNECTION_LIFETIME", "5m"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_CONNECTION_LIFETIME: %w", err)
	}

	connIdle, err := time.ParseDuration(getEnvOrDefault("DB_CONNECTION_IDLE", "5m"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_CONNECTION_IDLE: %w", err)
	}

	dbConfig := &DBConfig{
		Host:               getEnvOrDefault("DB_HOST", "localhost"),
		Port:               dbPort,
		User:               os.Getenv("DB_USER"),
		Password:           os.Getenv("DB_PASSWORD"),
		Name:               os.Getenv("DB_NAME"),
		SSLMode:            getEnvOrDefault("DB_SSLMODE", "disable"),
		MaxOpen:            maxOpen,
		MaxIdle:            maxIdle,
		ConnectionLifetime: connLifetime,
		ConnectionIdle:     connIdle,
	}

	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.SSLMode,
	)

	appPort := getEnvOrDefault("APP_PORT", "8080")
	// Ensure port has colon prefix for Fiber
	if appPort[0] != ':' {
		appPort = ":" + appPort
	}

	return &Config{
		AppPort:      appPort,
		DB:           dbConfig,
		DatabaseURL:  dbURL,
		GeminiAPIKey: os.Getenv("GEMINI_API_KEY"),
		ChromaDBURL:  getEnvOrDefault("CHROMADB_URL", "http://localhost:8000"),
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func ConnectDB(dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
