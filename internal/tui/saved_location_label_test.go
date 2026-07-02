package tui

import (
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestSavedLocationLabel(t *testing.T) {
	t.Parallel()

	location := domain.SavedLocation{
		Name:        "Madrid",
		Admin1:      "Community of Madrid",
		CountryCode: "ES",
		Latitude:    40.4168,
		Longitude:   -3.7038,
	}

	got := savedLocationLabel(location)
	expected := "Madrid, Community of Madrid, ES  40.4168, -3.7038"

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestSavedLocationLabelWithoutAdmin1(t *testing.T) {
	t.Parallel()

	location := domain.SavedLocation{
		Name:        "Copenhagen",
		CountryCode: "DK",
		Latitude:    55.6761,
		Longitude:   12.5683,
	}

	got := savedLocationLabel(location)
	expected := "Copenhagen, DK  55.6761, 12.5683"

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestSavedLocationLabelUnknownLocation(t *testing.T) {
	t.Parallel()

	got := savedLocationLabel(domain.SavedLocation{})
	expected := "Unknown location"

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
