package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port        int
	DatabaseURL string
	JWTSecret   string
	JWTExpiry   int64
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:        50051,
		DatabaseURL: getEnv("USER_DB_URL", "postgres://user:password@localhost:5432/sig_agro_users"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiry:   3600,
	}

	port := os.Getenv("USER_SERVICE_PORT")
	if port != "" {
		var err error
		cfg.Port, err = strconv.Atoi(port)
		if err != nil {
			log.Printf("Invalid port, using default: %d", cfg.Port)
		}
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
