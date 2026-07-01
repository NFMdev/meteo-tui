package tui

import (
	"fmt"
	"strings"

	"github.com/nfmdev/meteo/internal/domain"
)

func savedLocationLabel(location domain.SavedLocation) string {
	parts := []string{
		strings.TrimSpace(location.Name),
		strings.TrimSpace(location.Admin1),
		strings.TrimSpace(location.CountryCode),
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

	if location.Latitude == 0 && location.Longitude == 0 {
		return label
	}

	return fmt.Sprintf(
		"%s  %.4f, %.4f",
		label,
		location.Latitude,
		location.Longitude,
	)
}
