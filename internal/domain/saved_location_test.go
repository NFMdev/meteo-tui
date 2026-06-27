package domain

import "testing"

func TestSavedLocationFromSearchResult(t *testing.T) {
	t.Parallel()

	result := LocationSearchResult{
		Name:        " Copenhagen ",
		Country:     " Denmark ",
		CountryCode: " dk ",
		Admin1:      " Capital Region ",
		Latitude:    55.6761,
		Longitude:   12.5683,
		Timezone:    " Europe/Copenhagen ",
	}

	location := SavedLocationFromSearchResult(result)

	if location.Name != "Copenhagen" {
		t.Fatalf("expected name Copenhagen, got %q", location.Name)
	}

	if location.Country != "Denmark" {
		t.Fatalf("expected country Denmark, got %q", location.Country)
	}

	if location.CountryCode != "DK" {
		t.Fatalf("expected country code DK, got %q", location.CountryCode)
	}

	if location.Admin1 != "Capital Region" {
		t.Fatalf("expected admin1 Capital Region, got %q", location.Admin1)
	}

	if location.Latitude != 55.6761 {
		t.Fatalf("expected latitude 55.6761, got %.4f", location.Latitude)
	}

	if location.Longitude != 12.5683 {
		t.Fatalf("expected longitude 12.5683, got %.4f", location.Longitude)
	}

	if location.Timezone != "Europe/Copenhagen" {
		t.Fatalf("expected timezone Europe/Copenhagen, got %q", location.Timezone)
	}
}

func TestSavedLocationFromLocation(t *testing.T) {
	t.Parallel()

	location := Location{
		City:      " Copenhagen ",
		Country:   " dk ",
		Latitude:  55.6761,
		Longitude: 12.5683,
		Timezone:  " Europe/Copenhagen ",
	}

	saved := SavedLocationFromLocation(location)

	if saved.Name != "Copenhagen" {
		t.Fatalf("expected name Copenhagen, got %q", saved.Name)
	}

	if saved.CountryCode != "DK" {
		t.Fatalf("expected country code DK, got %q", saved.CountryCode)
	}

	if saved.Latitude != 55.6761 {
		t.Fatalf("expected latitude 55.6761, got %.4f", saved.Latitude)
	}

	if saved.Longitude != 12.5683 {
		t.Fatalf("expected longitude 12.5683, got %.4f", saved.Longitude)
	}

	if saved.Timezone != "Europe/Copenhagen" {
		t.Fatalf("expected timezone Europe/Copenhagen, got %q", saved.Timezone)
	}
}

func TestNormalizeSavedLocation(t *testing.T) {
	t.Parallel()

	location := NormalizeSavedLocation(SavedLocation{
		Name:        " Madrid ",
		Country:     " Spain ",
		CountryCode: " es ",
		Admin1:      " Community of Madrid ",
		Timezone:    " Europe/Madrid ",
	})

	if location.Name != "Madrid" {
		t.Fatalf("expected name Madrid, got %q", location.Name)
	}

	if location.Country != "Spain" {
		t.Fatalf("expected country Spain, got %q", location.Country)
	}

	if location.CountryCode != "ES" {
		t.Fatalf("expected country code ES, got %q", location.CountryCode)
	}

	if location.Admin1 != "Community of Madrid" {
		t.Fatalf("expected admin1 Community of Madrid, got %q", location.Admin1)
	}

	if location.Timezone != "Europe/Madrid" {
		t.Fatalf("expected timezone Europe/Madrid, got %q", location.Timezone)
	}
}

func TestSavedLocationKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		location SavedLocation
		expected string
	}{
		{
			name: "simple city",
			location: SavedLocation{
				Name:        "Copenhagen",
				CountryCode: "DK",
			},
			expected: "copenhagen_dk",
		},
		{
			name: "trims and normalizes",
			location: SavedLocation{
				Name:        "  Madrid  ",
				CountryCode: " es ",
			},
			expected: "madrid_es",
		},
		{
			name: "city with spaces",
			location: SavedLocation{
				Name:        "New York",
				CountryCode: "US",
			},
			expected: "new_york_us",
		},
		{
			name: "city with repeated separators",
			location: SavedLocation{
				Name:        "New---York",
				CountryCode: "US",
			},
			expected: "new_york_us",
		},
		{
			name: "city with unicode letters",
			location: SavedLocation{
				Name:        "São Paulo",
				CountryCode: "BR",
			},
			expected: "são_paulo_br",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := SavedLocationKey(tt.location)

			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestSavedLocationKeyReturnsEmptyForMissingName(t *testing.T) {
	t.Parallel()

	got := SavedLocationKey(SavedLocation{
		CountryCode: "DK",
	})

	if got != "" {
		t.Fatalf("expected empty key, got %q", got)
	}
}

func TestSavedLocationKeyReturnsEmptyForMissingCountryCode(t *testing.T) {
	t.Parallel()

	got := SavedLocationKey(SavedLocation{
		Name: "Copenhagen",
	})

	if got != "" {
		t.Fatalf("expected empty key, got %q", got)
	}
}

func TestSameSavedLocation(t *testing.T) {
	t.Parallel()

	left := SavedLocation{
		Name:        " Copenhagen ",
		CountryCode: " dk ",
	}

	right := SavedLocation{
		Name:        "Copenhagen",
		CountryCode: "DK",
	}

	if !SameSavedLocation(left, right) {
		t.Fatal("expected same saved location")
	}
}

func TestSameSavedLocationReturnsFalseForDifferentLocations(t *testing.T) {
	t.Parallel()

	left := SavedLocation{
		Name:        "Copenhagen",
		CountryCode: "DK",
	}

	right := SavedLocation{
		Name:        "Madrid",
		CountryCode: "ES",
	}

	if SameSavedLocation(left, right) {
		t.Fatal("expected different saved locations")
	}
}

func TestSameSavedLocationReturnsFalseForInvalidLocation(t *testing.T) {
	t.Parallel()

	left := SavedLocation{}
	right := SavedLocation{}

	if SameSavedLocation(left, right) {
		t.Fatal("expected invalid locations not to match")
	}
}
