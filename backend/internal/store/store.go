// Package store is the Postgres persistence layer (pgx).
package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/domain"
)

// ErrNotFound is returned when a row does not exist.
var ErrNotFound = errors.New("not found")

// LaunchStep is a persisted pipeline step.
type LaunchStep struct {
	ID          string          `json:"id"`
	BusinessID  string          `json:"business_id"`
	Seq         int             `json:"seq"`
	Type        domain.StepType `json:"step_type"`
	Provider    string          `json:"provider"`
	Title       string          `json:"title"`
	Mode        string          `json:"mode"`
	Status      string          `json:"status"`
	Request     json.RawMessage `json:"request"`
	Response    json.RawMessage `json:"response"`
	ExternalRef string          `json:"external_ref"`
	Error       string          `json:"error"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
}

// Store wraps a pgx connection pool.
type Store struct {
	pool *pgxpool.Pool
}

// New opens a connection pool to the given database URL.
func New(ctx context.Context, databaseURL string) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}
	cfg.MaxConns = 10
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}
	return &Store{pool: pool}, nil
}

// Close releases the pool.
func (s *Store) Close() { s.pool.Close() }

// CreateBusiness inserts a new launch application.
func (s *Store) CreateBusiness(ctx context.Context, b *domain.Business) error {
	b.ID = uuid.NewString()
	addr, _ := json.Marshal(b.Address)
	const q = `
		INSERT INTO businesses
			(id, country, entity_type, legal_name, founder_name, founder_email,
			 founder_phone, founder_id_number, address, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		RETURNING created_at, updated_at`
	return s.pool.QueryRow(ctx, q,
		b.ID, b.Country, b.EntityType, b.LegalName, b.FounderName, b.FounderEmail,
		b.FounderPhone, b.FounderIDNumber, addr, b.Status,
	).Scan(&b.CreatedAt, &b.UpdatedAt)
}

// GetBusiness fetches a single business by id.
func (s *Store) GetBusiness(ctx context.Context, id string) (*domain.Business, error) {
	const q = `
		SELECT id, country, entity_type, legal_name, founder_name, founder_email,
		       founder_phone, founder_id_number, address, status, created_at, updated_at
		FROM businesses WHERE id = $1`
	var b domain.Business
	var addr []byte
	err := s.pool.QueryRow(ctx, q, id).Scan(
		&b.ID, &b.Country, &b.EntityType, &b.LegalName, &b.FounderName, &b.FounderEmail,
		&b.FounderPhone, &b.FounderIDNumber, &addr, &b.Status, &b.CreatedAt, &b.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(addr, &b.Address)
	return &b, nil
}

// ListBusinesses returns the most recent launches.
func (s *Store) ListBusinesses(ctx context.Context, limit int) ([]domain.Business, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	const q = `
		SELECT id, country, entity_type, legal_name, founder_name, founder_email,
		       founder_phone, founder_id_number, address, status, created_at, updated_at
		FROM businesses ORDER BY created_at DESC LIMIT $1`
	rows, err := s.pool.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Business
	for rows.Next() {
		var b domain.Business
		var addr []byte
		if err := rows.Scan(
			&b.ID, &b.Country, &b.EntityType, &b.LegalName, &b.FounderName, &b.FounderEmail,
			&b.FounderPhone, &b.FounderIDNumber, &addr, &b.Status, &b.CreatedAt, &b.UpdatedAt,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(addr, &b.Address)
		out = append(out, b)
	}
	return out, rows.Err()
}

// UpdateBusinessStatus sets the overall status.
func (s *Store) UpdateBusinessStatus(ctx context.Context, id, status string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE businesses SET status=$2, updated_at=now() WHERE id=$1`, id, status)
	return err
}

// CreateSteps bulk-inserts the planned steps for a business in one transaction.
func (s *Store) CreateSteps(ctx context.Context, businessID string, plan []domain.PlannedStep) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // no-op after commit

	const q = `
		INSERT INTO launch_steps (id, business_id, seq, step_type, provider, title, mode, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,'pending')`
	for _, p := range plan {
		if _, err := tx.Exec(ctx, q,
			uuid.NewString(), businessID, p.Seq, p.Type, p.Provider, p.Title, p.Mode,
		); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

// GetSteps returns all steps for a business ordered by sequence.
func (s *Store) GetSteps(ctx context.Context, businessID string) ([]LaunchStep, error) {
	const q = `
		SELECT id, business_id, seq, step_type, provider, title, mode, status,
		       request, response, external_ref, error, created_at, updated_at, completed_at
		FROM launch_steps WHERE business_id=$1 ORDER BY seq ASC`
	rows, err := s.pool.Query(ctx, q, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSteps(rows)
}

// GetStep fetches a single step by id.
func (s *Store) GetStep(ctx context.Context, stepID string) (*LaunchStep, error) {
	const q = `
		SELECT id, business_id, seq, step_type, provider, title, mode, status,
		       request, response, external_ref, error, created_at, updated_at, completed_at
		FROM launch_steps WHERE id=$1`
	rows, err := s.pool.Query(ctx, q, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	steps, err := scanSteps(rows)
	if err != nil {
		return nil, err
	}
	if len(steps) == 0 {
		return nil, ErrNotFound
	}
	return &steps[0], nil
}

// NextPendingStep returns the lowest-seq step that has not completed yet.
func (s *Store) NextPendingStep(ctx context.Context, businessID string) (*LaunchStep, error) {
	const q = `
		SELECT id, business_id, seq, step_type, provider, title, mode, status,
		       request, response, external_ref, error, created_at, updated_at, completed_at
		FROM launch_steps
		WHERE business_id=$1 AND status IN ('pending','failed')
		ORDER BY seq ASC LIMIT 1`
	rows, err := s.pool.Query(ctx, q, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	steps, err := scanSteps(rows)
	if err != nil {
		return nil, err
	}
	if len(steps) == 0 {
		return nil, ErrNotFound
	}
	return &steps[0], nil
}

// MarkStepRunning flips a step into the running state.
func (s *Store) MarkStepRunning(ctx context.Context, stepID string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE launch_steps SET status='running', error='', updated_at=now() WHERE id=$1`, stepID)
	return err
}

// CompleteStep persists a successful step result.
func (s *Store) CompleteStep(ctx context.Context, stepID string, res domain.StepResult) error {
	data, _ := json.Marshal(res.Data)
	const q = `
		UPDATE launch_steps
		SET status='completed', response=$2, external_ref=$3, error='',
		    completed_at=now(), updated_at=now()
		WHERE id=$1`
	_, err := s.pool.Exec(ctx, q, stepID, data, res.ExternalRef)
	return err
}

// FailStep records a failure message against a step.
func (s *Store) FailStep(ctx context.Context, stepID, msg string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE launch_steps SET status='failed', error=$2, updated_at=now() WHERE id=$1`, stepID, msg)
	return err
}

func scanSteps(rows pgx.Rows) ([]LaunchStep, error) {
	var out []LaunchStep
	for rows.Next() {
		var st LaunchStep
		if err := rows.Scan(
			&st.ID, &st.BusinessID, &st.Seq, &st.Type, &st.Provider, &st.Title, &st.Mode,
			&st.Status, &st.Request, &st.Response, &st.ExternalRef, &st.Error,
			&st.CreatedAt, &st.UpdatedAt, &st.CompletedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, st)
	}
	return out, rows.Err()
}
