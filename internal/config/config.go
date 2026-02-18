package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	RedisAddr    string
	Concurrency  int
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SlackToken   string
	DiscordToken string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		Concurrency:  getEnvInt("WORKER_CONCURRENCY", 10),
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnvInt("SMTP_PORT", 587),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SlackToken:   getEnv("SLACK_TOKEN", ""),
		DiscordToken: getEnv("DISCORD_TOKEN", ""),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}
