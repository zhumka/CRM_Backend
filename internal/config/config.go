package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит все настройки приложения, загружаемые из окружения.
type Config struct {
	HTTPPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string
	JWTTTL    time.Duration

	GinMode string
}

// Load читает конфигурацию из переменных окружения.
// Файл .env подхватывается, если присутствует (для локальной разработки).
func Load() (*Config, error) {
	_ = godotenv.Load()

	ttlHours := getEnvInt("JWT_TTL_HOURS", 24)

	cfg := &Config{
		HTTPPort:   getEnv("HTTP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "crm"),
		DBPassword: getEnv("DB_PASSWORD", "crm"),
		DBName:     getEnv("DB_NAME", "crm"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
		JWTTTL:     time.Duration(ttlHours) * time.Hour,
		GinMode:    getEnv("GIN_MODE", "debug"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

// DSN возвращает строку подключения к PostgreSQL.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
