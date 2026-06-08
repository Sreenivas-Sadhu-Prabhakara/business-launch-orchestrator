// Command server is the business-launch-orchestrator HTTP API.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/api"
	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/auth"
	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/config"
	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/orchestrator"
	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/providers"
	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/store"
)

func main() {
	cfg := config.Load()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Connect to Postgres (retry briefly so `docker compose up` ordering is forgiving).
	st, err := connectWithRetry(ctx, cfg.DatabaseURL, cfg.DBMaxConns, 10, 2*time.Second)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer st.Close()

	if err := st.Migrate(ctx); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Println("migrations applied")

	authSvc := auth.New(cfg.JWTSecret, 24*time.Hour)
	if cfg.JWTSecret == "dev-insecure-secret-change-me" {
		log.Println("WARNING: using the default JWT secret — set JWT_SECRET for production")
	}
	seedUsers(ctx, st, cfg)

	reg := providers.NewRegistry(providers.Config{
		RazorpayKeyID:     cfg.RazorpayKeyID,
		RazorpayKeySecret: cfg.RazorpayKeySecret,
		StripeSecretKey:   cfg.StripeSecretKey,
		PayMongoSecretKey: cfg.PayMongoSecretKey,
		AnthropicAPIKey:   cfg.AnthropicAPIKey,
		AnthropicModel:    cfg.AnthropicModel,
		CSLAPIKey:         cfg.CSLAPIKey,
		ForceMock:         cfg.ForceMock,
	})
	svc := orchestrator.New(st, reg)
	handler := api.New(svc, st, authSvc, cfg.CORSOrigin)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("listening on :%s (force_mock=%v)", cfg.Port, cfg.ForceMock)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")
	shutdownCtx, shCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shCancel()
	_ = srv.Shutdown(shutdownCtx)
}

// seedUsers ensures the configured admin + demo accounts exist (idempotent).
func seedUsers(ctx context.Context, st *store.Store, cfg config.Config) {
	ensure := func(username, password, role string) {
		if username == "" || password == "" {
			return
		}
		hash, err := auth.HashPassword(password)
		if err != nil {
			log.Printf("seed %q: hash error: %v", username, err)
			return
		}
		created, err := st.EnsureUser(ctx, username, hash, role)
		if err != nil {
			log.Printf("seed %q: %v", username, err)
			return
		}
		if created {
			log.Printf("seeded %s account %q", role, username)
		}
	}
	ensure(cfg.AdminUsername, cfg.AdminPassword, auth.RoleAdmin)
	ensure(cfg.DemoUsername, cfg.DemoPassword, auth.RoleUser)
}

func connectWithRetry(ctx context.Context, url string, maxConns, attempts int, delay time.Duration) (*store.Store, error) {
	var lastErr error
	for i := 0; i < attempts; i++ {
		st, err := store.New(ctx, url, maxConns)
		if err == nil {
			return st, nil
		}
		lastErr = err
		log.Printf("database not ready (attempt %d/%d): %v", i+1, attempts, err)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}
	return nil, lastErr
}
