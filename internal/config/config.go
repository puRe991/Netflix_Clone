// Package config loads runtime configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	AppEnv               string
	Port                 string
	DatabaseURL          string
	JWTSecret            string
	AppURL               string
	StripeSecretKey      string
	StripeBasicPriceID   string
	StripePremiumPriceID string
	StripeWebhookSecret  string
}

func Load() (*Config, error) {
	loadDotEnv(".env")

	appEnv := getenv("APP_ENV", "development")

	cfg := &Config{
		AppEnv:               appEnv,
		Port:                 getenv("PORT", "3000"),
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		AppURL:               getenv("APP_URL", "http://localhost:3000"),
		StripeSecretKey:      os.Getenv("STRIPE_SECRET_KEY"),
		StripeBasicPriceID:   os.Getenv("STRIPE_BASIC_PRICE_ID"),
		StripePremiumPriceID: os.Getenv("STRIPE_PREMIUM_PRICE_ID"),
		StripeWebhookSecret:  os.Getenv("STRIPE_WEBHOOK_SECRET"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if cfg.JWTSecret == "" {
		if appEnv != "development" {
			return nil, fmt.Errorf("JWT_SECRET is required outside development")
		}
		cfg.JWTSecret = "dev-only-secret-change-me"
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// loadDotEnv reads KEY=VALUE pairs from a .env file (if present) into the
// process environment. Real environment variables always win: a var already
// set in the OS environment is never overwritten by the file, so deployments
// that inject config via the platform (Docker, systemd, ...) are unaffected.
// Missing files are not an error - .env is optional local convenience, and
// os.Getenv alone never reads it, which previously made DATABASE_URL/JWT_SECRET
// silently empty even though start.bat had created a valid .env.
func loadDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = unquote(strings.TrimSpace(value))
		if key == "" {
			continue
		}

		if _, set := os.LookupEnv(key); !set {
			os.Setenv(key, value)
		}
	}
}

func unquote(v string) string {
	if len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
			return v[1 : len(v)-1]
		}
	}
	return v
}
