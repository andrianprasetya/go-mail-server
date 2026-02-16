package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server
	AppPort string
	AppEnv  string

	// SMTP
	SMTPHost     string
	SMTPPort     int
	SMTPEmail    string
	SMTPPassword string

	// Email
	ReceiverEmail string

	// CORS
	AllowedOrigins []string

	// Rate Limiting
	RateLimit           int
	RateLimitExpiration int // in hours
}

// Configuration errors
var (
	ErrMissingSMTPEmail     = errors.New("SMTP_EMAIL is required")
	ErrMissingSMTPPassword  = errors.New("SMTP_PASSWORD is required")
	ErrMissingReceiverEmail = errors.New("RECEIVER_EMAIL is required")
)

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (development)
	if err := godotenv.Load(); err != nil {
		log.Println("[Config] No .env file found, using environment variables")
	}

	smtpPort, err := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	if err != nil {
		smtpPort = 587
	}

	rateLimit, err := strconv.Atoi(getEnv("RATE_LIMIT", "2"))
	if err != nil {
		rateLimit = 2
	}

	rateLimitExpiration, err := strconv.Atoi(getEnv("RATE_LIMIT_EXPIRATION_HOURS", "24"))
	if err != nil {
		rateLimitExpiration = 24
	}

	allowedOrigins := strings.Split(getEnv("ALLOWED_ORIGINS", "*"), ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	cfg := &Config{
		AppPort:             getEnv("APP_PORT", "3000"),
		AppEnv:              getEnv("APP_ENV", "development"),
		SMTPHost:            getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:            smtpPort,
		SMTPEmail:           getEnvWithFallback("SMTP_EMAIL", "SMTP_USERNAME", ""),
		SMTPPassword:        getEnv("SMTP_PASSWORD", ""),
		ReceiverEmail:       getEnv("RECEIVER_EMAIL", ""),
		AllowedOrigins:      allowedOrigins,
		RateLimit:           rateLimit,
		RateLimitExpiration: rateLimitExpiration,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if all required configuration is present
func (c *Config) Validate() error {
	if c.SMTPEmail == "" {
		return ErrMissingSMTPEmail
	}
	if c.SMTPPassword == "" {
		return ErrMissingSMTPPassword
	}
	if c.ReceiverEmail == "" {
		return ErrMissingReceiverEmail
	}
	return nil
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvWithFallback tries the primary key first, then falls back to an alternate key
func getEnvWithFallback(primary, fallback, defaultValue string) string {
	if value := os.Getenv(primary); value != "" {
		return value
	}
	if value := os.Getenv(fallback); value != "" {
		return value
	}
	return defaultValue
}
