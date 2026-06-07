package providers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/domain"
)

// paymentClients makes real sandbox calls to the supported payment gateways,
// falling back to a deterministic mock when keys are absent or FORCE_MOCK is set.
type paymentClients struct {
	cfg  Config
	http *http.Client
}

func (p *paymentClients) hasRazorpay() bool {
	return !p.cfg.ForceMock && p.cfg.RazorpayKeyID != "" && p.cfg.RazorpayKeySecret != ""
}
func (p *paymentClients) hasStripe() bool {
	return !p.cfg.ForceMock && p.cfg.StripeSecretKey != ""
}
func (p *paymentClients) hasPayMongo() bool {
	return !p.cfg.ForceMock && p.cfg.PayMongoSecretKey != ""
}

// Razorpay (India) — create a test-mode Order to prove live connectivity.
// Docs: https://razorpay.com/docs/api/orders/create/
func (p *paymentClients) Razorpay(ctx context.Context, b domain.Business) (domain.StepResult, error) {
	if !p.hasRazorpay() {
		return mockPayment("razorpay", "order", b), nil
	}
	form := url.Values{
		"amount":   {"100"}, // ₹1.00 in paise — a sandbox setup probe
		"currency": {"INR"},
		"receipt":  {"launch_" + b.ID},
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.razorpay.com/v1/orders", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(p.cfg.RazorpayKeyID, p.cfg.RazorpayKeySecret)

	body, err := p.do(req)
	if err != nil {
		return domain.StepResult{}, err
	}
	var out struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Amount int    `json:"amount"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return domain.StepResult{}, fmt.Errorf("razorpay decode: %w", err)
	}
	return domain.StepResult{
		ExternalRef: out.ID,
		Message:     "Razorpay test-mode merchant order created (live sandbox).",
		Data: map[string]any{
			"provider": "razorpay", "mode": "live", "order_id": out.ID,
			"status": out.Status, "raw": json.RawMessage(body),
		},
	}, nil
}

// Stripe (US) — create a test-mode Customer representing the merchant.
// Docs: https://docs.stripe.com/api/customers/create
func (p *paymentClients) Stripe(ctx context.Context, b domain.Business) (domain.StepResult, error) {
	if !p.hasStripe() {
		return mockPayment("stripe", "cus", b), nil
	}
	form := url.Values{
		"name":        {b.LegalName},
		"email":       {b.FounderEmail},
		"description": {"Merchant account for " + b.LegalName},
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.stripe.com/v1/customers", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+p.cfg.StripeSecretKey)

	body, err := p.do(req)
	if err != nil {
		return domain.StepResult{}, err
	}
	var out struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return domain.StepResult{}, fmt.Errorf("stripe decode: %w", err)
	}
	return domain.StepResult{
		ExternalRef: out.ID,
		Message:     "Stripe test-mode merchant customer created (live sandbox).",
		Data: map[string]any{
			"provider": "stripe", "mode": "live", "customer_id": out.ID,
			"email": out.Email, "raw": json.RawMessage(body),
		},
	}, nil
}

// PayMongo (Philippines) — create a test-mode Payment Link.
// Docs: https://developers.paymongo.com/reference/create-a-link
func (p *paymentClients) PayMongo(ctx context.Context, b domain.Business) (domain.StepResult, error) {
	if !p.hasPayMongo() {
		return mockPayment("paymongo", "link", b), nil
	}
	payload := map[string]any{
		"data": map[string]any{
			"attributes": map[string]any{
				"amount":      10000, // ₱100.00 in centavos
				"description": "Merchant setup for " + b.LegalName,
				"remarks":     "launch_" + b.ID,
			},
		},
	}
	raw, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.paymongo.com/v1/links", strings.NewReader(string(raw)))
	req.Header.Set("Content-Type", "application/json")
	// PayMongo uses HTTP basic auth: base64(secret_key + ":")
	auth := base64.StdEncoding.EncodeToString([]byte(p.cfg.PayMongoSecretKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)

	body, err := p.do(req)
	if err != nil {
		return domain.StepResult{}, err
	}
	var out struct {
		Data struct {
			ID         string `json:"id"`
			Attributes struct {
				CheckoutURL string `json:"checkout_url"`
				Status      string `json:"status"`
			} `json:"attributes"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return domain.StepResult{}, fmt.Errorf("paymongo decode: %w", err)
	}
	return domain.StepResult{
		ExternalRef: out.Data.ID,
		Message:     "PayMongo test-mode payment link created (live sandbox).",
		Data: map[string]any{
			"provider": "paymongo", "mode": "live", "link_id": out.Data.ID,
			"checkout_url": out.Data.Attributes.CheckoutURL,
			"status":       out.Data.Attributes.Status, "raw": json.RawMessage(body),
		},
	}, nil
}

// do executes the request and returns the body, treating >=400 as an error.
func (p *paymentClients) do(req *http.Request) ([]byte, error) {
	resp, err := p.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%s returned %d: %s", req.URL.Host, resp.StatusCode, string(body))
	}
	return body, nil
}

// mockPayment returns a deterministic merchant-setup result.
func mockPayment(provider, prefix string, b domain.Business) domain.StepResult {
	ref := fmt.Sprintf("%s_test_%s", prefix, digits(provider+b.ID, 14))
	return domain.StepResult{
		ExternalRef: ref,
		Message:     fmt.Sprintf("%s merchant configured (mock — set the API key to go live).", provider),
		Data: map[string]any{
			"provider": provider, "mode": "mock", "merchant_ref": ref,
			"currency": "local", "note": "Deterministic sandbox stand-in.",
		},
	}
}
