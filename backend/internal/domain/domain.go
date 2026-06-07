// Package domain holds the dependency-free core types shared across the
// orchestrator, the persistence layer and the provider adapters.
//
// Keeping these types in their own package avoids import cycles: both
// `store` (persistence) and `providers` (integrations) depend on `domain`,
// but never on each other.
package domain

import "time"

// Country is an ISO-3166 alpha-2 code for a jurisdiction we support.
type Country string

const (
	CountryIndia       Country = "IN"
	CountryPhilippines Country = "PH"
	CountryUS          Country = "US"
)

// Valid reports whether the country is one we have an adapter for.
func (c Country) Valid() bool {
	switch c {
	case CountryIndia, CountryPhilippines, CountryUS:
		return true
	default:
		return false
	}
}

// StepType is one logical stage in the "launch a business" pipeline. The same
// step type means different concrete API calls in each country (e.g.
// tax_registration is GST in India, IRS EIN in the US, BIR in the PH).
type StepType string

const (
	StepStrategyCheck    StepType = "strategy_check"
	StepFounderKYC       StepType = "founder_kyc"
	StepLiabilitiesCheck StepType = "liabilities_check"
	StepNameCheck        StepType = "name_check"
	StepIPCheck          StepType = "ip_check"
	StepEntityReg        StepType = "entity_registration"
	StepTaxReg           StepType = "tax_registration"
	StepRegistrations    StepType = "registrations"
	StepBankAccount      StepType = "bank_account"
	StepPaymentGateway   StepType = "payment_gateway"
	StepCompliance       StepType = "compliance_registration"
)

// Step / business status values.
const (
	StatusDraft      = "draft"
	StatusPending    = "pending"
	StatusRunning    = "running"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
	StatusInProgress = "in_progress"
)

// Integration modes surfaced in the plan so the UI can badge each step.
const (
	ModeLive = "live" // hits a real (sandbox) provider API
	ModeMock = "mock" // deterministic simulated response
)

// Address is the registered office / principal place of business.
type Address struct {
	Line1      string `json:"line1"`
	Line2      string `json:"line2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

// Business is a single end-to-end launch application.
type Business struct {
	ID              string    `json:"id"`
	Country         Country   `json:"country"`
	EntityType      string    `json:"entity_type"`
	LegalName       string    `json:"legal_name"`
	FounderName     string    `json:"founder_name"`
	FounderEmail    string    `json:"founder_email"`
	FounderPhone    string    `json:"founder_phone"`
	FounderIDNumber string    `json:"founder_id_number"`
	Address         Address   `json:"address"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// PlannedStep describes one step in a country's pipeline before it runs.
type PlannedStep struct {
	Seq      int      `json:"seq"`
	Type     StepType `json:"type"`
	Provider string   `json:"provider"`
	Title    string   `json:"title"`
	Mode     string   `json:"mode"`
}

// StepResult is what a provider adapter returns after executing a step.
type StepResult struct {
	ExternalRef string         `json:"external_ref"` // CIN, EIN, GSTIN, merchant id, ...
	Data        map[string]any `json:"data"`         // full structured response
	Message     string         `json:"message"`
}
