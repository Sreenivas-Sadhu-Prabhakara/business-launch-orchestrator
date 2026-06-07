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
	st, err := connectWithRetry(ctx, cfg.DatabaseURL, 10, 2*time.Second)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer st.Close()

	if err := st.Migrate(ctx); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Println("migrations applied")

	reg := providers.NewRegistry(providers.Config{
		RazorpayKeyID:     cfg.RazorpayKeyID,
		RazorpayKeySecret: cfg.RazorpayKeySecret,
		StripeSecretKey:   cfg.StripeSecretKey,
		PayMongoSecretKey: cfg.PayMongoSecretKey,
		ForceMock:         cfg.ForceMock,
	})
	svc := orchestrator.New(st, reg)
	handler := api.New(svc, st, cfg.CORSOrigin)

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

func connectWithRetry(ctx context.Context, url string, attempts int, delay time.Duration) (*store.Store, error) {
	var lastErr error
	for i := 0; i < attempts; i++ {
		st, err := store.New(ctx, url)
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
