package providers

import (
	"strings"
	"unicode"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/domain"
)

// countryDisplay returns a human-readable jurisdiction name.
func countryDisplay(c domain.Country) string {
	switch c {
	case domain.CountryIndia:
		return "India"
	case domain.CountryPhilippines:
		return "Philippines"
	case domain.CountryUS:
		return "United States"
	default:
		return string(c)
	}
}

// launchYear returns the incorporation year, defaulting to 2026 if unset.
func launchYear(b domain.Business) int {
	if y := b.CreatedAt.Year(); y >= 2000 {
		return y
	}
	return 2026
}

// firstLetter returns the first alphabetic character of s (uppercased), or "X".
func firstLetter(s string) string {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return strings.ToUpper(string(r))
		}
	}
	return "X"
}

// stateCode returns a 2-letter uppercase code derived from a state name, or def.
func stateCode(state, def string) string {
	letters := make([]rune, 0, 2)
	for _, r := range state {
		if unicode.IsLetter(r) {
			letters = append(letters, unicode.ToUpper(r))
			if len(letters) == 2 {
				break
			}
		}
	}
	if len(letters) == 2 {
		return string(letters)
	}
	return def
}
