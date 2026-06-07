package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/domain"
)

// usAdapter implements the US company-formation pipeline.
//
// Real upstreams (documented per step):
//   - KYC:        Persona / Middesk / Stripe Identity
//   - Name:       Secretary of State business name search
//   - Incorp:     Middesk Agents / Legalinc / Stripe Atlas (state LLC/Corp filing)
//   - Tax:        IRS EIN via Form SS-4
//   - Banking:    Mercury / Brex / Stripe Treasury
//   - Payments:   Stripe  ← LIVE sandbox
//   - Compliance: Registered agent + state tax / franchise registration
type usAdapter struct {
	cfg Config
	pay *paymentClients
}

func (a *usAdapter) Country() domain.Country { return domain.CountryUS }

func (a *usAdapter) Plan() []domain.PlannedStep {
	payMode := paymentMode(a.cfg, a.cfg.StripeSecretKey != "")
	return []domain.PlannedStep{
		{Seq: 1, Type: domain.StepFounderKYC, Provider: "Persona / Middesk", Title: "Verify founder identity", Mode: domain.ModeMock},
		{Seq: 2, Type: domain.StepNameCheck, Provider: "Secretary of State", Title: "Check entity name availability", Mode: domain.ModeMock},
		{Seq: 3, Type: domain.StepEntityReg, Provider: "Middesk Agents (state filing)", Title: "File Articles of Organization", Mode: domain.ModeMock},
		{Seq: 4, Type: domain.StepTaxReg, Provider: "IRS (Form SS-4)", Title: "Obtain Federal EIN", Mode: domain.ModeMock},
		{Seq: 5, Type: domain.StepBankAccount, Provider: "Mercury", Title: "Open business bank account", Mode: domain.ModeMock},
		{Seq: 6, Type: domain.StepPaymentGateway, Provider: "Stripe", Title: "Activate payment processing", Mode: payMode},
		{Seq: 7, Type: domain.StepCompliance, Provider: "Registered agent + state tax", Title: "Register agent & state tax", Mode: domain.ModeMock},
	}
}

func (a *usAdapter) Execute(ctx context.Context, step domain.StepType, b domain.Business) (domain.StepResult, error) {
	state := stateCode(b.Address.State, "DE")
	switch step {
	case domain.StepFounderKYC:
		return domain.StepResult{
			ExternalRef: "inq_" + digits("uskyc"+b.ID, 12),
			Message:     "Founder identity verified (government ID + selfie).",
			Data:        map[string]any{"status": "approved", "provider": "persona", "watchlist_hit": false},
		}, nil

	case domain.StepNameCheck:
		return domain.StepResult{
			ExternalRef: "name_" + digits("usname"+b.ID, 10),
			Message:     fmt.Sprintf("Name %q available in %s.", b.LegalName, state),
			Data:        map[string]any{"available": true, "state": state},
		}, nil

	case domain.StepEntityReg:
		filing := fmt.Sprintf("%s-%s-%s", state, strings.ToUpper(entityShort(b.EntityType)), digits("usfile"+b.ID, 7))
		return domain.StepResult{
			ExternalRef: filing,
			Message:     fmt.Sprintf("%s formed in %s; Articles filed.", b.EntityType, state),
			Data: map[string]any{
				"filing_number": filing, "state": state,
				"entity_type": b.EntityType, "document": "Articles of Organization",
			},
		}, nil

	case domain.StepTaxReg:
		ein := "88-" + digits("usein"+b.ID, 7)
		return domain.StepResult{
			ExternalRef: ein,
			Message:     "Federal EIN assigned by the IRS.",
			Data:        map[string]any{"ein": ein, "form": "SS-4"},
		}, nil

	case domain.StepBankAccount:
		acct := digits("usbank"+b.ID, 10)
		return domain.StepResult{
			ExternalRef: acct,
			Message:     "Business checking account opened with Mercury.",
			Data:        map[string]any{"account_number": acct, "routing_number": "084106768", "bank": "Mercury (Choice Financial)"},
		}, nil

	case domain.StepPaymentGateway:
		return a.pay.Stripe(ctx, b)

	case domain.StepCompliance:
		return domain.StepResult{
			ExternalRef: "ra_" + digits("usra"+b.ID, 10),
			Message:     "Registered agent appointed; state tax account opened.",
			Data: map[string]any{
				"registered_agent": "Middesk Agents, Inc.",
				"state_tax_id":     state + digits("ustax"+b.ID, 8),
			},
		}, nil
	}
	return domain.StepResult{}, fmt.Errorf("us: unsupported step %q", step)
}

// entityShort maps an entity type to the filing short code.
func entityShort(entityType string) string {
	t := strings.ToLower(entityType)
	switch {
	case strings.Contains(t, "c-corp"), strings.Contains(t, "c corp"), strings.Contains(t, "corp"):
		return "INC"
	case strings.Contains(t, "llc"):
		return "LLC"
	default:
		return "LLC"
	}
}
