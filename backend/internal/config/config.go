// Package config loads runtime configuration from environment variables.
package config

import (
	"os"
	"strconv"
)

// Config is the fully-resolved server configuration.
type Config struct {
	Port        string
	DatabaseURL string
	CORSOrigin  string

	// Payment provider sandbox credentials. When a key pair is empty the
	// corresponding payment step transparently falls back to a mock response
	// (the "hybrid" integration mode).
	RazorpayKeyID     string
	RazorpayKeySecret string
	StripeSecretKey   string
	PayMongoSecretKey string

	// ForceMock disables every live call (useful for offline demos / CI).
	ForceMock bool
}

// Load reads configuration from the environment, applying sane defaults.
func Load() Config {
	return Config{
		Port:        env("PORT", "8080"),
		DatabaseURL: env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/biz_launch?sslmode=disable"),
		CORSOrigin:  env("CORS_ORIGIN", "*"),

		RazorpayKeyID:     env("RAZORPAY_KEY_ID", ""),
		RazorpayKeySecret: env("RAZORPAY_KEY_SECRET", ""),
		StripeSecretKey:   env("STRIPE_SECRET_KEY", ""),
		PayMongoSecretKey: env("PAYMONGO_SECRET_KEY", ""),

		ForceMock: envBool("FORCE_MOCK", false),
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}
