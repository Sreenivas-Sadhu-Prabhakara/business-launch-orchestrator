// Package providers contains the per-country integration adapters that turn an
// abstract pipeline step (e.g. tax_registration) into concrete API calls.
//
// Integration strategy is "hybrid":
//   - Payment steps call REAL provider sandboxes (Razorpay / Stripe / PayMongo)
//     when test keys are configured, and fall back to a deterministic mock
//     otherwise.
//   - Government registry / KYC / banking steps are deterministic mocks whose
//     request+response shape mirrors the real upstream API, with the live
//     endpoint documented in code so a wiring is a drop-in replacement.
package providers

import (
	"context"
	"fmt"
	"hash/fnv"
	"net/http"
	"time"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/domain"
)

// Adapter executes the launch pipeline for one country.
type Adapter interface {
	Country() domain.Country
	Plan() []domain.PlannedStep
	Execute(ctx context.Context, step domain.StepType, b domain.Business) (domain.StepResult, error)
}

// Config carries provider credentials and toggles.
type Config struct {
	RazorpayKeyID     string
	RazorpayKeySecret string
	StripeSecretKey   string
	PayMongoSecretKey string
	AnthropicAPIKey   string
	AnthropicModel    string
	CSLAPIKey         string
	ForceMock         bool
}

// clients bundles the shared, country-agnostic integration clients that each
// adapter calls into.
type clients struct {
	cfg   Config
	pay   *paymentClients
	strat *strategyClient
	ip    *ipClient
	liab  *liabilitiesClient
}

// Registry maps a country to its adapter.
type Registry struct {
	adapters map[domain.Country]Adapter
}

// NewRegistry wires up the country adapters and the shared integration clients.
func NewRegistry(cfg Config) *Registry {
	httpc := &http.Client{Timeout: 25 * time.Second}
	cl := &clients{
		cfg:   cfg,
		pay:   &paymentClients{cfg: cfg, http: httpc},
		strat: &strategyClient{cfg: cfg, http: httpc},
		ip:    &ipClient{cfg: cfg, http: httpc},
		liab:  &liabilitiesClient{cfg: cfg, http: httpc},
	}
	return &Registry{
		adapters: map[domain.Country]Adapter{
			domain.CountryIndia:       &indiaAdapter{cfg: cfg, cl: cl},
			domain.CountryPhilippines: &phAdapter{cfg: cfg, cl: cl},
			domain.CountryUS:          &usAdapter{cfg: cfg, cl: cl},
		},
	}
}

// For returns the adapter for a country.
func (r *Registry) For(c domain.Country) (Adapter, bool) {
	a, ok := r.adapters[c]
	return a, ok
}

// Plan returns the pipeline plan for a country (empty if unsupported).
func (r *Registry) Plan(c domain.Country) []domain.PlannedStep {
	if a, ok := r.adapters[c]; ok {
		return a.Plan()
	}
	return nil
}

// ---- helpers shared by the mock adapters ----------------------------------

// seedNum derives a stable pseudo-random integer in [0, max) from a string, so
// generated identifiers (CIN, EIN, ...) are deterministic per business.
func seedNum(seed string, max uint64) uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(seed))
	return h.Sum64() % max
}

// digits returns an n-digit zero-padded number derived from seed.
func digits(seed string, n int) string {
	var max uint64 = 1
	for i := 0; i < n; i++ {
		max *= 10
	}
	return fmt.Sprintf("%0*d", n, seedNum(seed, max))
}

// modeFor reports whether a step will run live or mock: live only when a key
// (or keyless live capability) is present and FORCE_MOCK is off.
func modeFor(cfg Config, hasKey bool) string {
	if cfg.ForceMock || !hasKey {
		return domain.ModeMock
	}
	return domain.ModeLive
}
