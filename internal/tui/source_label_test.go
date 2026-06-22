package tui

import (
	"testing"
	"time"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestWeatherSourceLabelFreshOpenMeteo(t *testing.T) {
	t.Parallel()

	report := domain.WeatherReport{
		Source: domain.NewFreshWeatherSource(domain.WeatherProviderOpenMeteo),
	}

	got := weatherSourceLabel(report)

	if got != "Open-Meteo" {
		t.Fatalf("expected Open-Meteo, got %q", got)
	}
}

func TestWeatherSourceLabelCachedOpenMeteo(t *testing.T) {
	t.Parallel()

	cachedAt := time.Date(2026, 6, 21, 14, 30, 0, 0, time.UTC)

	report := domain.WeatherReport{
		Source: domain.NewCachedWeatherSource(
			domain.WeatherProviderOpenMeteo,
			cachedAt,
		),
	}

	got := weatherSourceLabel(report)
	expected := "Cached from Open-Meteo at 2026-06-21 14:30"

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestWeatherSourceLabelCachedWithoutTimestamp(t *testing.T) {
	t.Parallel()

	report := domain.WeatherReport{
		Source: domain.WeatherSource{
			Provider: domain.WeatherProviderOpenMeteo,
			Cached:   true,
		},
	}

	got := weatherSourceLabel(report)
	expected := "Cached from Open-Meteo"

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestWeatherSourceLabelUnknownFreshSource(t *testing.T) {
	t.Parallel()

	report := domain.WeatherReport{}

	got := weatherSourceLabel(report)

	if got != "Unknown source" {
		t.Fatalf("expected Unknown source, got %q", got)
	}
}

func TestFormatWeatherTimestamp(t *testing.T) {
	t.Parallel()

	value := time.Date(2026, 6, 21, 14, 30, 0, 0, time.UTC)

	got := formatWeatherTimestamp(value)

	if got != "2026-06-21 14:30" {
		t.Fatalf("expected 2026-06-21 14:30, got %q", got)
	}
}

func TestFormatWeatherTimestampZeroValue(t *testing.T) {
	t.Parallel()

	got := formatWeatherTimestamp(time.Time{})

	if got != "unknown time" {
		t.Fatalf("expected unknown time, got %q", got)
	}
}
