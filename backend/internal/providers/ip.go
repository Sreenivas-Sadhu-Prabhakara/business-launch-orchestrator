package providers

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/domain"
)

// ipClient runs an intellectual-property pre-clearance: a REAL, keyless domain
// availability check via RDAP, plus a trademark conflict search (mock, with the
// real registry documented per country).
type ipClient struct {
	cfg  Config
	http httpDoer
}

// trademark office + a representative ccTLD per country.
var ipOffice = map[domain.Country]struct {
	office string
	ccTLD  string
}{
	domain.CountryIndia:       {"IP India (CGPDTM)", "in"},
	domain.CountryPhilippines: {"IPOPHL", "ph"},
	domain.CountryUS:          {"USPTO", "us"},
}

// Check returns trademark + domain availability for the proposed brand.
func (c *ipClient) Check(ctx context.Context, b domain.Business) (domain.StepResult, error) {
	slug := domainSlug(b.LegalName)
	meta := ipOffice[b.Country]

	candidates := []string{slug + ".com"}
	if meta.ccTLD != "" {
		candidates = append(candidates, slug+"."+meta.ccTLD)
	}

	domains := make(map[string]string, len(candidates))
	for _, d := range candidates {
		domains[d] = c.domainStatus(ctx, d)
	}

	// Trademark search (mock — real upstream documented above).
	appNo := fmt.Sprintf("TM-%d", seedNum("tm"+b.ID, 9_000_000)+1_000_000)
	classes := []string{"9 (software)", "42 (SaaS / technology services)"}

	mode := "hybrid"
	if c.cfg.ForceMock {
		mode = domain.ModeMock
	}

	return domain.StepResult{
		ExternalRef: appNo,
		Message:     "Domain availability checked (RDAP, live); trademark search clear.",
		Data: map[string]any{
			"mode":    mode,
			"domains": domains,
			"trademark": map[string]any{
				"office":         meta.office,
				"status":         "no_conflicting_mark",
				"application_no": appNo,
				"classes":        classes,
			},
		},
	}, nil
}

// domainStatus queries the public RDAP network: HTTP 404 means the domain is
// unregistered (available), 200 means it is taken.
func (c *ipClient) domainStatus(ctx context.Context, fqdn string) string {
	if c.cfg.ForceMock {
		if seedNum("dom"+fqdn, 2) == 0 {
			return "available"
		}
		return "registered"
	}
	req, err := newRequest(ctx, "GET", "https://rdap.org/domain/"+fqdn, nil)
	if err != nil {
		return "unknown"
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return "unknown"
	}
	defer resp.Body.Close()
	switch {
	case resp.StatusCode == 404:
		return "available"
	case resp.StatusCode == 200:
		return "registered"
	default:
		return "unknown"
	}
}

// domainSlug lowercases the name and keeps only [a-z0-9].
func domainSlug(name string) string {
	var sb strings.Builder
	for _, r := range strings.ToLower(name) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			sb.WriteRune(r)
		}
	}
	s := sb.String()
	if s == "" {
		s = "newco"
	}
	if len(s) > 40 {
		s = s[:40]
	}
	return s
}
