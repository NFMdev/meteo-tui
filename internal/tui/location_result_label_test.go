package tui

import (
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestLocationSearchResultLabel(t *testing.T) {
	t.Parallel()

	result := domain.LocationSearchResult{
		Name:        "Copenhagen",
		Admin1:      "Capital Region",
		CountryCode: "DK",
		Latitude:    55.6761,
		Longitude:   12.5683,
	}

	got := locationSearchResultLabel(result)
	expected := "Copenhagen, Capital Region, DK  55.6761, 12.5683"

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestLocationSearchResultLabelWithoutAdmin1(t *testing.T) {
	t.Parallel()

	result := domain.LocationSearchResult{
		Name:        "Madrid",
		CountryCode: "ES",
		Latitude:    40.4168,
		Longitude:   -3.7038,
	}

	got := locationSearchResultLabel(result)
	expected := "Madrid, ES  40.4168, -3.7038"

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestLocationSearchResultLabelUnknownLocation(t *testing.T) {
	t.Parallel()

	got := locationSearchResultLabel(domain.LocationSearchResult{})
	expected := "Unknown location"

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
