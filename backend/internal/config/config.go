package config

import "os"

// Config holds the application configuration values sourced from the environment.
type Config struct {
	DatabaseURL string
	JWTSecret   string
	ServerPort  string
}

// Load reads environment variables and populates a Config struct with sane defaults for local development.
func Load() Config {
	return Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/mockapi?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key"),
		ServerPort:  getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
