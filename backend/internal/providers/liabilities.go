package providers

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/domain"
)

// liabilitiesClient screens the founder/entity for risk: sanctions/PEP exposure
// plus country-specific liability checks (director disqualification, litigation,
// liens, tax-defaulter status). Sanctions screening is REAL via the trade.gov
// Consolidated Screening List when CSL_API_KEY is set; everything else is a
// documented mock.
type liabilitiesClient struct {
	cfg  Config
	http httpDoer
}

// Screen runs the diligence checks and returns an overall risk verdict.
func (c *liabilitiesClient) Screen(ctx context.Context, b domain.Business) (domain.StepResult, error) {
	sanctions := c.sanctions(ctx, b)
	checks := countryLiabilityChecks(b.Country)

	risk := "low"
	if matches, _ := sanctions["matches"].(int); matches > 0 {
		risk = "review"
	}

	return domain.StepResult{
		ExternalRef: "risk:" + risk,
		Message:     "Liabilities & sanctions screening complete.",
		Data: map[string]any{
			"risk_level": risk,
			"sanctions":  sanctions,
			"checks":     checks,
		},
	}, nil
}

// sanctions hits the trade.gov Consolidated Screening List API when keyed.
func (c *liabilitiesClient) sanctions(ctx context.Context, b domain.Business) map[string]any {
	if c.cfg.ForceMock || c.cfg.CSLAPIKey == "" {
		return map[string]any{"status": "clear", "source": "mock", "matches": 0}
	}
	q := url.Values{"name": {b.FounderName}, "api_key": {c.cfg.CSLAPIKey}}
	endpoint := "https://api.trade.gov/consolidated_screening_list/search?" + q.Encode()

	req, err := newRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return map[string]any{"status": "error", "source": "trade.gov CSL", "matches": 0}
	}
	body, err := doJSON(c.http, req)
	if err != nil {
		return map[string]any{"status": "error", "source": "trade.gov CSL", "matches": 0, "error": err.Error()}
	}
	var out struct {
		Total int `json:"total"`
	}
	_ = json.Unmarshal(body, &out)
	status := "clear"
	if out.Total > 0 {
		status = "review"
	}
	return map[string]any{"status": status, "source": "trade.gov CSL (live)", "matches": out.Total}
}

// countryLiabilityChecks returns the mock per-jurisdiction diligence items.
func countryLiabilityChecks(c domain.Country) map[string]string {
	switch c {
	case domain.CountryIndia:
		return map[string]string{
			"MCA director disqualification (Sec 164)": "clear",
			"GST defaulter list":                      "clear",
			"CIBIL commercial bureau":                 "satisfactory",
			"eCourts pending litigation":              "none_found",
		}
	case domain.CountryPhilippines:
		return map[string]string{
			"SEC good-standing":           "clear",
			"BIR open-case / delinquency": "clear",
			"Court litigation search":     "none_found",
			"Credit (CIC) standing":       "satisfactory",
		}
	case domain.CountryUS:
		return map[string]string{
			"OFAC SDN list":           "clear",
			"UCC lien search":         "none",
			"PACER litigation search": "none_found",
			"Bankruptcy filings":      "none",
		}
	default:
		return map[string]string{"general_screening": "clear"}
	}
}
