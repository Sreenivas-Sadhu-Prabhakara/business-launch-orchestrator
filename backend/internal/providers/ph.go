package providers

import (
	"context"
	"fmt"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/domain"
)

// phAdapter implements the Philippine business-registration pipeline.
//
// Real upstreams (documented per step):
//   - KYC:        HyperVerge / Jumio (PhilID / government ID verification)
//   - Name:       SEC name verification (or DTI BNRS for sole proprietorship)
//   - Register:   SEC eSPARC / OneSEC (corporations) or DTI BNRS
//   - Tax:        BIR registration → TIN + Certificate of Registration (Form 2303)
//   - Banking:    UnionBank / BPI business account
//   - Payments:   PayMongo  ← LIVE sandbox
//   - Compliance: SSS + PhilHealth + Pag-IBIG employer registration + Mayor's permit
type phAdapter struct {
	cfg Config
	cl  *clients
}

func (a *phAdapter) Country() domain.Country { return domain.CountryPhilippines }

func (a *phAdapter) Plan() []domain.PlannedStep {
	payMode := modeFor(a.cfg, a.cfg.PayMongoSecretKey != "")
	return []domain.PlannedStep{
		{Seq: 1, Type: domain.StepStrategyCheck, Provider: "Claude (AI strategist)", Title: "Strategy & viability assessment", Mode: modeFor(a.cfg, a.cfg.AnthropicAPIKey != "")},
		{Seq: 2, Type: domain.StepFounderKYC, Provider: "HyperVerge (PhilID)", Title: "Verify founder identity", Mode: domain.ModeMock},
		{Seq: 3, Type: domain.StepLiabilitiesCheck, Provider: "trade.gov CSL + SEC/BIR/CIC", Title: "Liabilities & sanctions screening", Mode: modeFor(a.cfg, a.cfg.CSLAPIKey != "")},
		{Seq: 4, Type: domain.StepNameCheck, Provider: "SEC name verification", Title: "Verify company name", Mode: domain.ModeMock},
		{Seq: 5, Type: domain.StepIPCheck, Provider: "RDAP + IPOPHL", Title: "Trademark & domain check", Mode: modeFor(a.cfg, true)},
		{Seq: 6, Type: domain.StepEntityReg, Provider: "SEC eSPARC / OneSEC", Title: "Register corporation (SEC)", Mode: domain.ModeMock},
		{Seq: 7, Type: domain.StepTaxReg, Provider: "BIR (Form 1903/2303)", Title: "Register with BIR (TIN)", Mode: domain.ModeMock},
		{Seq: 8, Type: domain.StepRegistrations, Provider: "Mayor's permit + Barangay + DTI", Title: "Permits & registrations", Mode: domain.ModeMock},
		{Seq: 9, Type: domain.StepBankAccount, Provider: "UnionBank", Title: "Open corporate bank account", Mode: domain.ModeMock},
		{Seq: 10, Type: domain.StepPaymentGateway, Provider: "PayMongo", Title: "Activate payment gateway", Mode: payMode},
		{Seq: 11, Type: domain.StepCompliance, Provider: "SSS + PhilHealth + Pag-IBIG", Title: "Register as employer", Mode: domain.ModeMock},
	}
}

func (a *phAdapter) Execute(ctx context.Context, step domain.StepType, b domain.Business) (domain.StepResult, error) {
	year := launchYear(b)
	switch step {
	case domain.StepStrategyCheck:
		return a.cl.strat.Assess(ctx, b)

	case domain.StepLiabilitiesCheck:
		return a.cl.liab.Screen(ctx, b)

	case domain.StepIPCheck:
		return a.cl.ip.Check(ctx, b)

	case domain.StepFounderKYC:
		return domain.StepResult{
			ExternalRef: "kyc_" + digits("phkyc"+b.ID, 12),
			Message:     "Founder PhilID verified.",
			Data:        map[string]any{"status": "verified", "id_type": "PhilID", "liveness": "pass"},
		}, nil

	case domain.StepNameCheck:
		return domain.StepResult{
			ExternalRef: "NV-" + digits("phname"+b.ID, 8),
			Message:     fmt.Sprintf("Company name %q verified with SEC.", b.LegalName),
			Data:        map[string]any{"available": true},
		}, nil

	case domain.StepEntityReg:
		sec := fmt.Sprintf("CS%d%s", year, digits("phsec"+b.ID, 6))
		return domain.StepResult{
			ExternalRef: sec,
			Message:     "Corporation registered with the SEC.",
			Data: map[string]any{
				"sec_registration_number": sec, "system": "eSPARC / OneSEC",
				"document": "Certificate of Incorporation",
			},
		}, nil

	case domain.StepTaxReg:
		tin := fmt.Sprintf("%s-%s-%s-000", digits("phtin1"+b.ID, 3), digits("phtin2"+b.ID, 3), digits("phtin3"+b.ID, 3))
		return domain.StepResult{
			ExternalRef: tin,
			Message:     "Registered with the BIR; Certificate of Registration (2303) issued.",
			Data:        map[string]any{"tin": tin, "form": "2303", "rdo_code": digits("phrdo"+b.ID, 3)},
		}, nil

	case domain.StepRegistrations:
		permit := "BP-" + digits("phbp"+b.ID, 9)
		return domain.StepResult{
			ExternalRef: permit,
			Message:     "Local permits & registrations filed.",
			Data: map[string]any{
				"mayors_business_permit": permit,
				"barangay_clearance":     "BC-" + digits("phbc"+b.ID, 8),
				"dti_bnrs":               "BN-" + digits("phdti"+b.ID, 8),
				"fda_ltb":                "LTO-" + digits("phfda"+b.ID, 9),
			},
		}, nil

	case domain.StepBankAccount:
		acct := digits("phbank"+b.ID, 12)
		return domain.StepResult{
			ExternalRef: acct,
			Message:     "Corporate account opened with UnionBank.",
			Data:        map[string]any{"account_number": acct, "bank": "UnionBank of the Philippines"},
		}, nil

	case domain.StepPaymentGateway:
		return a.cl.pay.PayMongo(ctx, b)

	case domain.StepCompliance:
		return domain.StepResult{
			ExternalRef: "SSS-" + digits("phsss"+b.ID, 10),
			Message:     "Registered as employer with SSS, PhilHealth and Pag-IBIG.",
			Data: map[string]any{
				"sss_employer_no":        digits("phsss"+b.ID, 10),
				"philhealth_employer_no": digits("phphil"+b.ID, 12),
				"pagibig_employer_id":    digits("phpag"+b.ID, 12),
			},
		}, nil
	}
	return domain.StepResult{}, fmt.Errorf("ph: unsupported step %q", step)
}
