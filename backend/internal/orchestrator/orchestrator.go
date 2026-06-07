// Package orchestrator drives the end-to-end launch pipeline: it persists the
// plan, runs steps in order against the country adapters, and tracks state so a
// launch can be advanced one step at a time or run to completion.
package orchestrator

import (
	"context"
	"errors"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/domain"
	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/providers"
	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/store"
)

var (
	// ErrUnsupportedCountry is returned for a country with no adapter.
	ErrUnsupportedCountry = errors.New("unsupported country")
	// ErrNoPendingSteps means the pipeline has already completed.
	ErrNoPendingSteps = errors.New("no pending steps")
)

// Service is the orchestration use-case layer.
type Service struct {
	store *store.Store
	reg   *providers.Registry
}

// New constructs an orchestrator service.
func New(s *store.Store, r *providers.Registry) *Service {
	return &Service{store: s, reg: r}
}

// Registry exposes the provider registry (used by the API to render plans).
func (svc *Service) Registry() *providers.Registry { return svc.reg }

// CreateLaunch validates input, persists the business, and seeds its step plan.
func (svc *Service) CreateLaunch(ctx context.Context, b *domain.Business) error {
	if !b.Country.Valid() {
		return ErrUnsupportedCountry
	}
	if _, ok := svc.reg.For(b.Country); !ok {
		return ErrUnsupportedCountry
	}
	b.Status = domain.StatusDraft
	if err := svc.store.CreateBusiness(ctx, b); err != nil {
		return err
	}
	plan := svc.reg.Plan(b.Country)
	return svc.store.CreateSteps(ctx, b.ID, plan)
}

// AdvanceOne runs the next pending (or previously failed) step.
func (svc *Service) AdvanceOne(ctx context.Context, businessID string) (*store.LaunchStep, error) {
	b, err := svc.store.GetBusiness(ctx, businessID)
	if err != nil {
		return nil, err
	}
	next, err := svc.store.NextPendingStep(ctx, businessID)
	if errors.Is(err, store.ErrNotFound) {
		_ = svc.store.UpdateBusinessStatus(ctx, businessID, domain.StatusCompleted)
		return nil, ErrNoPendingSteps
	}
	if err != nil {
		return nil, err
	}
	return svc.runStep(ctx, b, next)
}

// RunAll executes every remaining step, stopping at the first failure.
func (svc *Service) RunAll(ctx context.Context, businessID string) ([]store.LaunchStep, error) {
	var ran []store.LaunchStep
	for {
		step, err := svc.AdvanceOne(ctx, businessID)
		if errors.Is(err, ErrNoPendingSteps) {
			return ran, nil
		}
		if err != nil {
			return ran, err
		}
		ran = append(ran, *step)
		if step.Status == domain.StatusFailed {
			return ran, nil
		}
	}
}

// runStep executes a single step against the country adapter and persists the
// outcome, returning the refreshed step row.
func (svc *Service) runStep(ctx context.Context, b *domain.Business, step *store.LaunchStep) (*store.LaunchStep, error) {
	adapter, ok := svc.reg.For(b.Country)
	if !ok {
		return nil, ErrUnsupportedCountry
	}

	_ = svc.store.UpdateBusinessStatus(ctx, b.ID, domain.StatusInProgress)
	if err := svc.store.MarkStepRunning(ctx, step.ID); err != nil {
		return nil, err
	}

	res, execErr := adapter.Execute(ctx, step.Type, *b)
	if execErr != nil {
		if err := svc.store.FailStep(ctx, step.ID, execErr.Error()); err != nil {
			return nil, err
		}
		return svc.store.GetStep(ctx, step.ID)
	}

	if err := svc.store.CompleteStep(ctx, step.ID, res); err != nil {
		return nil, err
	}

	// If that was the last step, flip the business to completed.
	if _, err := svc.store.NextPendingStep(ctx, b.ID); errors.Is(err, store.ErrNotFound) {
		_ = svc.store.UpdateBusinessStatus(ctx, b.ID, domain.StatusCompleted)
	}

	return svc.store.GetStep(ctx, step.ID)
}
