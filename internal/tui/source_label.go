package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/nfmdev/meteo/internal/domain"
)

func weatherSourceLabel(report domain.WeatherReport) string {
	source := report.Source

	provider := strings.TrimSpace(source.Provider)
	if provider == "" || provider == domain.WeatherProviderUnknown {
		provider = "unknown provider"
	}

	if source.Cached {
		if source.CachedAt.IsZero() {
			return fmt.Sprintf("Cached from %s", provider)
		}

		return fmt.Sprintf("Cached from %s at %s", provider, formatWeatherTimestamp(source.CachedAt))
	}

	if provider == "unknown provider" {
		return "Unknown source"
	}

	return provider
}

func formatWeatherTimestamp(value time.Time) string {
	if value.IsZero() {
		return "unknown time"
	}
	return value.Format("2006-01-02 15:04")
}
