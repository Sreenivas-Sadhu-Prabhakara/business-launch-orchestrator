package providers

import (
	"context"
	"testing"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/domain"
)

func sampleBusiness(c domain.Country) domain.Business {
	return domain.Business{
		ID:           "test-" + string(c),
		Country:      c,
		EntityType:   "Test Entity",
		LegalName:    "Acme Test Labs",
		FounderName:  "Jane Doe",
		FounderEmail: "jane@example.com",
		Address:      domain.Address{State: "Karnataka", Country: string(c)},
	}
}

// TestCountryPipelines runs every step of every country end-to-end in mock mode
// and asserts each produces a non-empty reference with no error.
func TestCountryPipelines(t *testing.T) {
	reg := NewRegistry(Config{ForceMock: true})
	for _, c := range []domain.Country{domain.CountryIndia, domain.CountryPhilippines, domain.CountryUS} {
		adapter, ok := reg.For(c)
		if !ok {
			t.Fatalf("%s: no adapter registered", c)
		}
		plan := adapter.Plan()
		if len(plan) != 7 {
			t.Errorf("%s: expected 7 planned steps, got %d", c, len(plan))
		}
		b := sampleBusiness(c)
		for _, p := range plan {
			res, err := adapter.Execute(context.Background(), p.Type, b)
			if err != nil {
				t.Errorf("%s/%s: unexpected error: %v", c, p.Type, err)
				continue
			}
			if res.ExternalRef == "" {
				t.Errorf("%s/%s: empty external_ref", c, p.Type)
			}
			if p.Type == domain.StepPaymentGateway && p.Mode != domain.ModeMock {
				t.Errorf("%s: payment step should be mock when FORCE_MOCK is set", c)
			}
		}
	}
}

// TestDeterministicRefs ensures generated identifiers are stable per business.
func TestDeterministicRefs(t *testing.T) {
	reg := NewRegistry(Config{ForceMock: true})
	a, _ := reg.For(domain.CountryIndia)
	b := sampleBusiness(domain.CountryIndia)

	r1, _ := a.Execute(context.Background(), domain.StepEntityReg, b)
	r2, _ := a.Execute(context.Background(), domain.StepEntityReg, b)
	if r1.ExternalRef != r2.ExternalRef {
		t.Errorf("CIN not deterministic: %q vs %q", r1.ExternalRef, r2.ExternalRef)
	}
	if r1.ExternalRef == "" {
		t.Error("expected a non-empty CIN")
	}
}

// TestUnsupportedStep verifies adapters reject unknown step types.
func TestUnsupportedStep(t *testing.T) {
	reg := NewRegistry(Config{ForceMock: true})
	a, _ := reg.For(domain.CountryUS)
	if _, err := a.Execute(context.Background(), domain.StepType("bogus"), sampleBusiness(domain.CountryUS)); err == nil {
		t.Error("expected error for unsupported step type")
	}
}
