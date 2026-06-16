package config

import "os"

type Config struct {
	DBPath     string
	SessionKey string
	Port       string
	Debug      bool
	AppEnv     string
}

func Load() *Config {
	return &Config{
		DBPath:     getEnv("DB_PATH", "./data/schoolbooks.db"),
		SessionKey: getEnv("SESSION_KEY", "change-in-production"),
		Port:       getEnv("PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),
		Debug:      false,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
