package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	JWTSecret    string
	JWTExpiresIn string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppPort:      getEnv("APP_PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "lebah"),
		DBPassword:   getEnv("DB_PASSWORD", "ganteng"),
		DBName:       getEnv("DB_NAME", "familymart"),
		DBSSLMode:    getEnv("DB_SSLMODE", "disable"),
		JWTSecret:    getEnv("JWT_SECRET", ""),
		JWTExpiresIn: getEnv("JWT_EXPIRES_IN", "15m"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
