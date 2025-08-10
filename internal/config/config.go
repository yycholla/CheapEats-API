package config

import (
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	API      APIConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type APIConfig struct {
	GooglePlacesAPIKey string
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "cheapeats"),
			Password: getEnv("DB_PASSWORD", "cheapeats_pass"),
			DBName:   getEnv("DB_NAME", "cheapeats_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		API: APIConfig{
			GooglePlacesAPIKey: getEnv("GOOGLE_PLACES_API_KEY", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}