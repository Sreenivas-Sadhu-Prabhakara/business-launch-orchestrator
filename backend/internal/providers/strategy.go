package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/domain"
)

// strategyClient produces a go/no-go business strategy assessment. When an
// Anthropic API key is configured it calls Claude (with a cached system prompt);
// otherwise it returns a deterministic mock assessment.
type strategyClient struct {
	cfg  Config
	http httpDoer
}

func (s *strategyClient) live() bool { return !s.cfg.ForceMock && s.cfg.AnthropicAPIKey != "" }

// assessment is the structured strategy output we ask Claude to return.
type assessment struct {
	RecommendedEntity string   `json:"recommended_entity"`
	ViabilityScore    int      `json:"viability_score"`
	KeyRisks          []string `json:"key_risks"`
	GoToMarket        []string `json:"go_to_market"`
	Summary           string   `json:"summary"`
}

// strategySystemPrompt is large and static so it can be prompt-cached across
// requests (Anthropic `cache_control: ephemeral`).
const strategySystemPrompt = `You are a senior startup-formation strategist advising founders on launching a new company. You weigh jurisdiction, entity structure, regulatory burden, tax exposure, market timing and capital needs.

Given a venture, produce a concise, pragmatic assessment. Be specific to the jurisdiction named. Consider:
- whether the proposed entity type is optimal (vs alternatives in that country)
- the top regulatory / tax / compliance risks for that jurisdiction
- 3-5 concrete go-to-market moves for the first 90 days
- an overall viability score from 0-100

Respond with STRICT JSON only (no prose, no markdown fences) matching exactly:
{"recommended_entity": string, "viability_score": number, "key_risks": string[], "go_to_market": string[], "summary": string}`

// Assess returns a strategy assessment for the venture.
func (s *strategyClient) Assess(ctx context.Context, b domain.Business) (domain.StepResult, error) {
	if !s.live() {
		return mockStrategy(b), nil
	}

	model := s.cfg.AnthropicModel
	if model == "" {
		model = "claude-sonnet-4-6"
	}

	userPrompt := fmt.Sprintf(
		"Venture to assess:\n- Country: %s\n- Proposed entity type: %s\n- Company name: %s\n- Founder: %s\n- Registered region: %s\n\nReturn the JSON assessment.",
		countryDisplay(b.Country), b.EntityType, b.LegalName, b.FounderName, b.Address.State,
	)

	reqBody := map[string]any{
		"model":      model,
		"max_tokens": 1024,
		"system": []map[string]any{
			{
				"type":          "text",
				"text":          strategySystemPrompt,
				"cache_control": map[string]string{"type": "ephemeral"},
			},
		},
		"messages": []map[string]any{
			{"role": "user", "content": userPrompt},
		},
	}
	raw, _ := json.Marshal(reqBody)

	req, _ := newRequest(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(raw))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", s.cfg.AnthropicAPIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	body, err := doJSON(s.http, req)
	if err != nil {
		return domain.StepResult{}, fmt.Errorf("anthropic: %w", err)
	}

	var out struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage map[string]any `json:"usage"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return domain.StepResult{}, fmt.Errorf("anthropic decode: %w", err)
	}

	var text strings.Builder
	for _, c := range out.Content {
		if c.Type == "text" {
			text.WriteString(c.Text)
		}
	}

	a := parseAssessment(text.String(), b)
	return domain.StepResult{
		ExternalRef: fmt.Sprintf("viability:%d", a.ViabilityScore),
		Message:     "AI strategy assessment generated (Claude, live).",
		Data: map[string]any{
			"mode": "live", "model": model, "assessment": a, "usage": out.Usage,
		},
	}, nil
}

// parseAssessment extracts the JSON object from Claude's reply, tolerating any
// surrounding text; falls back to a seeded assessment on failure.
func parseAssessment(text string, b domain.Business) assessment {
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start >= 0 && end > start {
		var a assessment
		if err := json.Unmarshal([]byte(text[start:end+1]), &a); err == nil && a.ViabilityScore > 0 {
			return a
		}
	}
	a := mockAssessment(b)
	if strings.TrimSpace(text) != "" {
		a.Summary = strings.TrimSpace(text)
	}
	return a
}

func mockStrategy(b domain.Business) domain.StepResult {
	a := mockAssessment(b)
	return domain.StepResult{
		ExternalRef: fmt.Sprintf("viability:%d", a.ViabilityScore),
		Message:     "Strategy assessment generated (mock — set ANTHROPIC_API_KEY for a live Claude analysis).",
		Data:        map[string]any{"mode": "mock", "assessment": a},
	}
}

func mockAssessment(b domain.Business) assessment {
	score := 62 + int(seedNum("strategy"+b.ID, 33)) // 62..94
	risks := map[domain.Country][]string{
		domain.CountryIndia: {
			"GST compliance cadence and input-credit reconciliation",
			"FDI/FEMA reporting if foreign capital is raised",
			"State-level labour and Shops & Establishment variance",
		},
		domain.CountryPhilippines: {
			"SEC + BIR + LGU permit sequencing can delay go-live",
			"Foreign equity caps under the Foreign Investments Act",
			"Withholding-tax and books-of-account formalities",
		},
		domain.CountryUS: {
			"State-by-state nexus for sales tax once you sell across lines",
			"FinCEN Beneficial Ownership (BOI) reporting deadlines",
			"Franchise tax / annual report obligations by state",
		},
	}[b.Country]
	return assessment{
		RecommendedEntity: b.EntityType,
		ViabilityScore:    score,
		KeyRisks:          risks,
		GoToMarket: []string{
			"Validate pricing with 10 design-partner conversations before launch",
			"Stand up payments + a single high-intent acquisition channel first",
			"Instrument activation metrics from day one",
		},
		Summary: fmt.Sprintf(
			"A %s in %s is a sound default for this stage; structure is appropriate and the regulatory path is well-trodden. Prioritise distribution over breadth in the first 90 days.",
			b.EntityType, countryDisplay(b.Country),
		),
	}
}
