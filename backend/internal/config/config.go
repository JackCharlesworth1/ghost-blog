package config

import (
	"os"
)

type Config struct {
	DatabaseURL  string
	Port         string
	AdminPassword string
	Environment  string
}

func Load() *Config {
	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ghostblog?sslmode=disable"),
		Port:          getEnv("PORT", "8080"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
