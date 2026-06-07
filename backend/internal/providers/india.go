package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/domain"
)

// indiaAdapter implements the Indian incorporation pipeline.
//
// Real upstreams (documented per step):
//   - KYC:        SurePass / Karza / Signzy (PAN + Aadhaar verification)
//   - Name:       MCA "RUN" (Reserve Unique Name)
//   - Incorp:     MCA SPICe+ (INC-32) → Certificate of Incorporation + CIN
//   - Tax:        Income Tax (PAN/TAN) + GSTN (GSTIN), bundled into SPICe+ AGILE-PRO
//   - Banking:    RazorpayX / Open current account
//   - Payments:   Razorpay  ← LIVE sandbox
//   - Compliance: EPFO + ESIC registration (also via SPICe+ AGILE-PRO)
type indiaAdapter struct {
	cfg Config
	pay *paymentClients
}

func (a *indiaAdapter) Country() domain.Country { return domain.CountryIndia }

func (a *indiaAdapter) Plan() []domain.PlannedStep {
	payMode := paymentMode(a.cfg, a.cfg.RazorpayKeyID != "" && a.cfg.RazorpayKeySecret != "")
	return []domain.PlannedStep{
		{Seq: 1, Type: domain.StepFounderKYC, Provider: "SurePass (PAN + Aadhaar)", Title: "Verify founder identity", Mode: domain.ModeMock},
		{Seq: 2, Type: domain.StepNameCheck, Provider: "MCA RUN", Title: "Reserve company name", Mode: domain.ModeMock},
		{Seq: 3, Type: domain.StepEntityReg, Provider: "MCA SPICe+ (INC-32)", Title: "Incorporate company (CIN)", Mode: domain.ModeMock},
		{Seq: 4, Type: domain.StepTaxReg, Provider: "Income Tax + GSTN", Title: "Obtain PAN, TAN & GSTIN", Mode: domain.ModeMock},
		{Seq: 5, Type: domain.StepBankAccount, Provider: "RazorpayX", Title: "Open current account", Mode: domain.ModeMock},
		{Seq: 6, Type: domain.StepPaymentGateway, Provider: "Razorpay", Title: "Activate payment gateway", Mode: payMode},
		{Seq: 7, Type: domain.StepCompliance, Provider: "EPFO + ESIC", Title: "Register for PF & ESI", Mode: domain.ModeMock},
	}
}

func (a *indiaAdapter) Execute(ctx context.Context, step domain.StepType, b domain.Business) (domain.StepResult, error) {
	year := launchYear(b)
	switch step {
	case domain.StepFounderKYC:
		return domain.StepResult{
			ExternalRef: "kyc_" + digits("inkyc"+b.ID, 12),
			Message:     "Founder PAN & Aadhaar verified.",
			Data: map[string]any{
				"pan_status": "VALID", "aadhaar_status": "VERIFIED",
				"name_match": true, "provider": "surepass",
			},
		}, nil

	case domain.StepNameCheck:
		return domain.StepResult{
			ExternalRef: "RUN" + digits("inname"+b.ID, 10),
			Message:     fmt.Sprintf("Name %q reserved with MCA.", b.LegalName),
			Data:        map[string]any{"available": true, "reserved_until_days": 20},
		}, nil

	case domain.StepEntityReg:
		cin := indiaCIN(b, year)
		return domain.StepResult{
			ExternalRef: cin,
			Message:     "Certificate of Incorporation issued via SPICe+.",
			Data: map[string]any{
				"cin": cin, "form": "INC-32 (SPICe+)",
				"roc": "RoC-Bangalore", "incorporation_year": year,
			},
		}, nil

	case domain.StepTaxReg:
		pan := indiaPAN(b)
		gstin := "29" + pan + "1Z" + strings.ToUpper(digits("ingst"+b.ID, 1))
		tan := "BLR" + strings.ToUpper(firstLetter(b.LegalName)) + digits("intan"+b.ID, 5) + "K"
		return domain.StepResult{
			ExternalRef: gstin,
			Message:     "PAN, TAN and GSTIN allotted.",
			Data:        map[string]any{"pan": pan, "tan": tan, "gstin": gstin},
		}, nil

	case domain.StepBankAccount:
		acct := digits("inbank"+b.ID, 12)
		return domain.StepResult{
			ExternalRef: acct,
			Message:     "Current account opened with RazorpayX.",
			Data:        map[string]any{"account_number": acct, "ifsc": "RATN0VAAPIS", "bank": "RazorpayX (RBL)"},
		}, nil

	case domain.StepPaymentGateway:
		return a.pay.Razorpay(ctx, b)

	case domain.StepCompliance:
		return domain.StepResult{
			ExternalRef: "PF-" + digits("inpf"+b.ID, 10),
			Message:     "Registered with EPFO (PF) and ESIC (ESI).",
			Data: map[string]any{
				"epfo_code": "KNBNG" + digits("inepf"+b.ID, 7),
				"esic_code": "53" + digits("inesi"+b.ID, 15),
			},
		}, nil
	}
	return domain.StepResult{}, fmt.Errorf("india: unsupported step %q", step)
}

func indiaCIN(b domain.Business, year int) string {
	// U<5-digit industry><2-char state>YYYY PTC <6-digit>
	state := stateCode(b.Address.State, "KA")
	return fmt.Sprintf("U72900%s%dPTC%s", state, year, digits("incin"+b.ID, 6))
}

func indiaPAN(b domain.Business) string {
	// 5 letters (4th = 'C' for company) + 4 digits + 1 letter.
	return "AABC" + strings.ToUpper(firstLetter(b.LegalName)) + digits("inpan"+b.ID, 4) + "Z"
}
