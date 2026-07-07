// Package config loads runtime configuration from environment variables.
package config

import (
	"fmt"
	"os"
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
