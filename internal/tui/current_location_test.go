package tui

import (
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestCurrentSavedLocationReturnsCurrentWeatherLocation(t *testing.T) {
	t.Parallel()

	model := Model{
		report: domain.WeatherReport{
			Location: domain.Location{
				City:      "Copenhagen",
				Country:   "DK",
				Latitude:  55.6761,
				Longitude: 12.5683,
				Timezone:  "Europe/Copenhagen",
			},
		},
	}

	location, ok := model.currentSavedLocation()
	if !ok {
		t.Fatal("expected current saved location")
	}

	if location.Name != "Copenhagen" {
		t.Fatalf("expected Copenhagen, got %q", location.Name)
	}

	if location.CountryCode != "DK" {
		t.Fatalf("expected DK, got %q", location.CountryCode)
	}

	if location.Latitude != 55.6761 {
		t.Fatalf("expected latitude 55.6761, got %.4f", location.Latitude)
	}

	if location.Longitude != 12.5683 {
		t.Fatalf("expected longitude 12.5683, got %.4f", location.Longitude)
	}

	if location.Timezone != "Europe/Copenhagen" {
		t.Fatalf("expected Europe/Copenhagen, got %q", location.Timezone)
	}
}

func TestCurrentSavedLocationReturnsFalseWhenLocationIsMissing(t *testing.T) {
	t.Parallel()

	model := Model{}

	_, ok := model.currentSavedLocation()
	if ok {
		t.Fatal("expected no current saved location")
	}
}
