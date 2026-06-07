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
	DBMaxConns  int

	// Payment provider sandbox credentials. When a key pair is empty the
	// corresponding payment step transparently falls back to a mock response
	// (the "hybrid" integration mode).
	RazorpayKeyID     string
	RazorpayKeySecret string
	StripeSecretKey   string
	PayMongoSecretKey string

	// AnthropicAPIKey powers the live AI strategy assessment (Claude). When
	// empty the strategy step returns a deterministic mock assessment.
	AnthropicAPIKey string
	AnthropicModel  string

	// CSLAPIKey enables live sanctions/liabilities screening via the
	// trade.gov Consolidated Screening List API.
	CSLAPIKey string

	// ForceMock disables every live call (useful for offline demos / CI).
	ForceMock bool
}

// Load reads configuration from the environment, applying sane defaults.
func Load() Config {
	return Config{
		Port:        env("PORT", "8080"),
		DatabaseURL: env("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/biz_launch?sslmode=disable"),
		CORSOrigin:  env("CORS_ORIGIN", "*"),
		DBMaxConns:  envInt("DB_MAX_CONNS", 10),

		RazorpayKeyID:     env("RAZORPAY_KEY_ID", ""),
		RazorpayKeySecret: env("RAZORPAY_KEY_SECRET", ""),
		StripeSecretKey:   env("STRIPE_SECRET_KEY", ""),
		PayMongoSecretKey: env("PAYMONGO_SECRET_KEY", ""),

		AnthropicAPIKey: env("ANTHROPIC_API_KEY", ""),
		AnthropicModel:  env("ANTHROPIC_MODEL", "claude-sonnet-4-6"),
		CSLAPIKey:       env("CSL_API_KEY", ""),

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

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}
