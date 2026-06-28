package config

import (
	"errors"
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
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
		DefaultCity:    "Copenhagen",
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
				DefaultCity:    "Copenhagen",
				DefaultCountry: tt.country,
			})

			if !errors.Is(err, ErrInvalidCountryCode) {
				t.Fatalf("expected ErrInvalidCountryCode, got %v", err)
			}
		})
	}
}

func TestNormalizeConfigNormalizesFavorites(t *testing.T) {
	t.Parallel()

	config := NormalizeConfig(AppConfig{
		DefaultCity:    " Copenhagen ",
		DefaultCountry: " dk ",
		Favorites: []domain.SavedLocation{
			{
				Name:        " Madrid ",
				Country:     " Spain ",
				CountryCode: " es ",
				Admin1:      " Community of Madrid ",
				Timezone:    " Europe/Madrid ",
			},
		},
	})

	if config.DefaultCity != "Copenhagen" {
		t.Fatalf("expected default city Copenhagen, got %q", config.DefaultCity)
	}

	if config.DefaultCountry != "DK" {
		t.Fatalf("expected default country DK, got %q", config.DefaultCountry)
	}

	if len(config.Favorites) != 1 {
		t.Fatalf("expected 1 favorite, got %d", len(config.Favorites))
	}

	favorite := config.Favorites[0]

	if favorite.Name != "Madrid" {
		t.Fatalf("expected favorite name Madrid, got %q", favorite.Name)
	}

	if favorite.Country != "Spain" {
		t.Fatalf("expected favorite country Spain, got %q", favorite.Country)
	}

	if favorite.CountryCode != "ES" {
		t.Fatalf("expected favorite country code ES, got %q", favorite.CountryCode)
	}

	if favorite.Admin1 != "Community of Madrid" {
		t.Fatalf("expected favorite admin1 Community of Madrid, got %q", favorite.Admin1)
	}

	if favorite.Timezone != "Europe/Madrid" {
		t.Fatalf("expected favorite timezone Europe/Madrid, got %q", favorite.Timezone)
	}
}

func TestValidateConfigAllowsFavorites(t *testing.T) {
	t.Parallel()

	config := AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "DK",
		Favorites: []domain.SavedLocation{
			{
				Name:        "Madrid",
				CountryCode: "ES",
			},
			{
				Name:        "Copenhagen",
				CountryCode: "DK",
			},
		},
	}

	err := ValidateConfig(config)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateConfigRejectsFavoriteWithoutName(t *testing.T) {
	t.Parallel()

	config := AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "DK",
		Favorites: []domain.SavedLocation{
			{
				CountryCode: "ES",
			},
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrFavoriteNameRequired) {
		t.Fatalf("expected ErrFavoriteNameRequired, got %v", err)
	}
}

func TestValidateConfigRejectsFavoriteWithoutCountryCode(t *testing.T) {
	t.Parallel()

	config := AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "DK",
		Favorites: []domain.SavedLocation{
			{
				Name: "Madrid",
			},
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrFavoriteCountryCodeRequired) {
		t.Fatalf("expected ErrFavoriteCountryCodeRequired, got %v", err)
	}
}

func TestValidateConfigRejectsFavoriteWithInvalidCountryCode(t *testing.T) {
	t.Parallel()

	config := AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "DK",
		Favorites: []domain.SavedLocation{
			{
				Name:        "Madrid",
				CountryCode: "ESP",
			},
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrInvalidCountryCode) {
		t.Fatalf("expected ErrInvalidCountryCode, got %v", err)
	}
}

func TestValidateConfigRejectsDuplicateFavorites(t *testing.T) {
	t.Parallel()

	config := AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "DK",
		Favorites: []domain.SavedLocation{
			{
				Name:        " Madrid ",
				CountryCode: " es ",
			},
			{
				Name:        "Madrid",
				CountryCode: "ES",
			},
		},
	}

	err := ValidateConfig(config)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrDuplicateFavorite) {
		t.Fatalf("expected ErrDuplicateFavorite, got %v", err)
	}
}
