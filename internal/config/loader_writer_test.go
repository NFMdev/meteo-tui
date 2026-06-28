package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestLoadConfigLoadsValidConfig(t *testing.T) {
	t.Parallel()

	path := writeTestConfigFile(t, `{
		"default_city": "Copenhagen",
		"default_country": "DK"
	}`)

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if config.DefaultCity != "Copenhagen" {
		t.Fatalf("expected city Copenhagen, got %q", config.DefaultCity)
	}

	if config.DefaultCountry != "DK" {
		t.Fatalf("expected country DK, got %q", config.DefaultCountry)
	}
}

func TestLoadConfigNormalizesConfig(t *testing.T) {
	t.Parallel()

	path := writeTestConfigFile(t, `{
		"default_city": "  Copenhagen  ",
		"default_country": " dk "
	}`)

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if config.DefaultCity != "Copenhagen" {
		t.Fatalf("expected normalized city Copenhagen, got %q", config.DefaultCity)
	}

	if config.DefaultCountry != "DK" {
		t.Fatalf("expected normalized country DK, got %q", config.DefaultCountry)
	}
}

func TestLoadConfigReturnsConfigNotFound(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "missing.json")

	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("expected ErrConfigNotFound, got %v", err)
	}
}

func TestLoadConfigRejectsMalformedJSON(t *testing.T) {
	t.Parallel()

	path := writeTestConfigFile(t, `{ invalid json`)

	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "decode config file") {
		t.Fatalf("expected decode config error, got %v", err)
	}
}

func TestLoadConfigRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	path := writeTestConfigFile(t, `{
		"default_city": "Copenhagen",
		"default_contry": "DK"
	}`)

	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("expected unknown field error, got %v", err)
	}
}

func TestLoadConfigRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	path := writeTestConfigFile(t, `{
		"default_city": "",
		"default_country": "DK"
	}`)

	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrDefaultCityRequired) {
		t.Fatalf("expected ErrDefaultCityRequired, got %v", err)
	}
}

func TestWriteConfigCreatesParentDirectoryAndWritesConfig(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "nested", "meteo", "config.json")

	err := WriteConfig(path, AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "dk",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("expected written config to load, got %v", err)
	}

	if config.DefaultCity != "Copenhagen" {
		t.Fatalf("expected city Copenhagen, got %q", config.DefaultCity)
	}

	if config.DefaultCountry != "DK" {
		t.Fatalf("expected country DK, got %q", config.DefaultCountry)
	}
}

func TestWriteConfigRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")

	err := WriteConfig(path, AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "DNK",
	})

	if !errors.Is(err, ErrInvalidCountryCode) {
		t.Fatalf("expected ErrInvalidCountryCode, got %v", err)
	}
}

func TestWriteConfigRejectsEmptyPath(t *testing.T) {
	t.Parallel()

	err := WriteConfig("", AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "DK",
	})

	if !errors.Is(err, ErrConfigPathRequired) {
		t.Fatalf("expected ErrConfigPathRequired, got %v", err)
	}
}

func TestLoadConfigRejectsEmptyPath(t *testing.T) {
	t.Parallel()

	_, err := LoadConfig("")

	if !errors.Is(err, ErrConfigPathRequired) {
		t.Fatalf("expected ErrConfigPathRequired, got %v", err)
	}
}

func TestLoadConfigAllowsMissingFavoritesField(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")

	content := `{
		"default_city": "Copenhagen",
		"default_country": "DK"
	}`

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if config.DefaultCity != "Copenhagen" {
		t.Fatalf("expected default city Copenhagen, got %q", config.DefaultCity)
	}

	if config.DefaultCountry != "DK" {
		t.Fatalf("expected default country DK, got %q", config.DefaultCountry)
	}

	if len(config.Favorites) != 0 {
		t.Fatalf("expected no favorites, got %d", len(config.Favorites))
	}
}

func TestLoadConfigLoadsFavorites(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")

	content := `{
		"default_city": "Copenhagen",
		"default_country": "DK",
		"favorites": [
			{
				"name": " Madrid ",
				"country": " Spain ",
				"country_code": " es ",
				"admin1": " Community of Madrid ",
				"latitude": 40.4168,
				"longitude": -3.7038,
				"timezone": " Europe/Madrid "
			}
		]
	}`

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	config, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(config.Favorites) != 1 {
		t.Fatalf("expected 1 favorite, got %d", len(config.Favorites))
	}

	favorite := config.Favorites[0]

	if favorite.Name != "Madrid" {
		t.Fatalf("expected Madrid, got %q", favorite.Name)
	}

	if favorite.Country != "Spain" {
		t.Fatalf("expected Spain, got %q", favorite.Country)
	}

	if favorite.CountryCode != "ES" {
		t.Fatalf("expected ES, got %q", favorite.CountryCode)
	}

	if favorite.Admin1 != "Community of Madrid" {
		t.Fatalf("expected Community of Madrid, got %q", favorite.Admin1)
	}

	if favorite.Latitude != 40.4168 {
		t.Fatalf("expected latitude 40.4168, got %.4f", favorite.Latitude)
	}

	if favorite.Longitude != -3.7038 {
		t.Fatalf("expected longitude -3.7038, got %.4f", favorite.Longitude)
	}

	if favorite.Timezone != "Europe/Madrid" {
		t.Fatalf("expected Europe/Madrid, got %q", favorite.Timezone)
	}
}

func TestLoadConfigRejectsDuplicateFavorites(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")

	content := `{
		"default_city": "Copenhagen",
		"default_country": "DK",
		"favorites": [
			{
				"name": "Madrid",
				"country_code": "ES"
			},
			{
				"name": " Madrid ",
				"country_code": " es "
			}
		]
	}`

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrDuplicateFavorite) {
		t.Fatalf("expected ErrDuplicateFavorite, got %v", err)
	}
}

func TestWriteConfigWritesFavorites(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")

	config := AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "DK",
		Favorites: []domain.SavedLocation{
			{
				Name:        "Madrid",
				Country:     "Spain",
				CountryCode: "ES",
				Admin1:      "Community of Madrid",
				Latitude:    40.4168,
				Longitude:   -3.7038,
				Timezone:    "Europe/Madrid",
			},
		},
	}

	err := WriteConfig(path, config)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	loaded, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("expected no error loading written config, got %v", err)
	}

	if len(loaded.Favorites) != 1 {
		t.Fatalf("expected 1 favorite, got %d", len(loaded.Favorites))
	}

	if loaded.Favorites[0].Name != "Madrid" {
		t.Fatalf("expected Madrid, got %q", loaded.Favorites[0].Name)
	}

	if loaded.Favorites[0].CountryCode != "ES" {
		t.Fatalf("expected ES, got %q", loaded.Favorites[0].CountryCode)
	}
}

func writeTestConfigFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.json")

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	return path
}
