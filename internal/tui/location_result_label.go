package tui

import (
	"fmt"
	"strings"

	"github.com/nfmdev/meteo/internal/domain"
)

func locationSearchResultLabel(result domain.LocationSearchResult) string {
	parts := []string{
		strings.TrimSpace(result.Name),
		strings.TrimSpace(result.Admin1),
		strings.TrimSpace(result.CountryCode),
	}

	visibleParts := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			visibleParts = append(visibleParts, part)
		}
	}

	label := strings.Join(visibleParts, ", ")
	if label == "" {
		label = "Unknown location"
	}

	if result.Latitude == 0 && result.Longitude == 0 {
		return label
	}

	return fmt.Sprintf(
		"%s  %.4f, %.4f",
		label,
		result.Latitude,
		result.Longitude,
	)
}
