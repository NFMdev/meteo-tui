package config

import (
	"errors"
	"testing"
)

func TestNormalizeConfig(t *testing.T) {
	t.Parallel()

	config := NormalizeConfig(AppConfig{
		DefaultCity:    "  Copenhagen  ",
		DefaultCountry: "  dk  ",
	})

	if config.DefaultCity != "Copenhagen" {
		t.Fatalf("expected city Copenhagen, got %q", config.DefaultCity)
	}

	if config.DefaultCountry != "DK" {
		t.Fatalf("expected country DK, got %q", config.DefaultCountry)
	}
}

func TestValidationConfigAcceptsValidConfig(t *testing.T) {
	t.Parallel()

	err := ValidateConfig(AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "DK",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateConfigRejectsEmptyCity(t *testing.T) {
	t.Parallel()

	err := ValidateConfig(AppConfig{
		DefaultCity:    "",
		DefaultCountry: "DK",
	})

	if !errors.Is(err, ErrDefaultCityRequired) {
		t.Fatalf("expected ErrDefaultCityRequired, got %v", err)
	}
}

func TestValidateConfigRejectsEmptyCountry(t *testing.T) {
	t.Parallel()

	err := ValidateConfig(AppConfig{
		DefaultCity:    "Aalborg",
		DefaultCountry: "",
	})

	if !errors.Is(err, ErrDefaultCountryRequired) {
		t.Fatalf("expected ErrDefaultCountryRequired, got %v", err)
	}
}

func TestValidateConfigRejectsInvalidCountryCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		country string
	}{
		{
			name:    "one letter",
			country: "D",
		},
		{
			name:    "three letters",
			country: "DNK",
		},
		{
			name:    "contains number",
			country: "D1",
		},
		{
			name:    "contains symbol",
			country: "D-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(AppConfig{
				DefaultCity:    "Aalborg",
				DefaultCountry: tt.country,
			})

			if !errors.Is(err, ErrInvalidCountryCode) {
				t.Fatalf("expected ErrInvalidCountryCode, got %v", err)
			}
		})
	}
}
