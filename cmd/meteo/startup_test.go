package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	meteoConfig "github.com/nfmdev/meteo/internal/config"
)

func TestResolveStartupOptionsUsesFallbackWhenNoConfigAndNoFlags(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "missing-config.json")

	resolved, err := resolveStartupOptions(startupOptions{
		configPath: configPath,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resolved.city != fallbackCity {
		t.Fatalf("expected fallback city %q, got %q", fallbackCity, resolved.city)
	}

	if resolved.country != fallbackCountry {
		t.Fatalf("expected fallback country %q, got %q", fallbackCountry, resolved.country)
	}

	if resolved.configPath != configPath {
		t.Fatalf("expected config path %q, got %q", configPath, resolved.configPath)
	}
}

func TestResolveStartupOptionsUsesConfigWhenNoFlagsProvided(t *testing.T) {
	t.Parallel()

	configPath := writeStartupTestConfig(t, meteoConfig.AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "DK",
	})

	resolved, err := resolveStartupOptions(startupOptions{
		configPath: configPath,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resolved.city != "Copenhagen" {
		t.Fatalf("expected city Copenhagen, got %q", resolved.city)
	}

	if resolved.country != "DK" {
		t.Fatalf("expected country DK, got %q", resolved.country)
	}
}

func TestResolveStartupOptionsCliFlagsOverrideConfig(t *testing.T) {
	t.Parallel()

	configPath := writeStartupTestConfig(t, meteoConfig.AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "DK",
	})

	resolved, err := resolveStartupOptions(startupOptions{
		city:       "Madrid",
		country:    "ES",
		configPath: configPath,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resolved.city != "Madrid" {
		t.Fatalf("expected city Madrid, got %q", resolved.city)
	}

	if resolved.country != "ES" {
		t.Fatalf("expected country ES, got %q", resolved.country)
	}
}

func TestResolveStartupOptionsAllowsPartialCityOverride(t *testing.T) {
	t.Parallel()

	configPath := writeStartupTestConfig(t, meteoConfig.AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "DK",
	})

	resolved, err := resolveStartupOptions(startupOptions{
		city:       "Madrid",
		configPath: configPath,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resolved.city != "Madrid" {
		t.Fatalf("expected city Madrid, got %q", resolved.city)
	}

	if resolved.country != "DK" {
		t.Fatalf("expected country DK from config, got %q", resolved.country)
	}
}

func TestResolveStartupOptionsAllowsPartialCountryOverride(t *testing.T) {
	t.Parallel()

	configPath := writeStartupTestConfig(t, meteoConfig.AppConfig{
		DefaultCity:    "Copenhagen",
		DefaultCountry: "DK",
	})

	resolved, err := resolveStartupOptions(startupOptions{
		country:    "ES",
		configPath: configPath,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resolved.city != "Copenhagen" {
		t.Fatalf("expected city Copenhagen from config, got %q", resolved.city)
	}

	if resolved.country != "ES" {
		t.Fatalf("expected country ES, got %q", resolved.country)
	}
}

func TestResolveStartupOptionsNormalizesCountry(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "missing-config.json")

	resolved, err := resolveStartupOptions(startupOptions{
		city:       "Copenhagen",
		country:    " dk ",
		configPath: configPath,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resolved.country != "DK" {
		t.Fatalf("expected country DK, got %q", resolved.country)
	}
}

func TestResolveStartupOptionsTrimsCity(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "missing-config.json")

	resolved, err := resolveStartupOptions(startupOptions{
		city:       "  Copenhagen  ",
		country:    "DK",
		configPath: configPath,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if resolved.city != "Copenhagen" {
		t.Fatalf("expected city Copenhagen, got %q", resolved.city)
	}
}

func TestResolveStartupOptionsReturnsErrorForInvalidCountry(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "missing-config.json")

	_, err := resolveStartupOptions(startupOptions{
		city:       "Copenhagen",
		country:    "DNK",
		configPath: configPath,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, meteoConfig.ErrInvalidCountryCode) {
		t.Fatalf("expected ErrInvalidCountryCode, got %v", err)
	}
}

func TestResolveStartupOptionsReturnsErrorForInvalidConfig(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "config.json")

	if err := os.WriteFile(configPath, []byte(`{
		"default_city": "",
		"default_country": "DK"
	}`), 0o644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := resolveStartupOptions(startupOptions{
		configPath: configPath,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, meteoConfig.ErrDefaultCityRequired) {
		t.Fatalf("expected ErrDefaultCityRequired, got %v", err)
	}
}

func TestResolveStartupOptionsReturnsErrorForMalformedConfig(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "config.json")

	if err := os.WriteFile(configPath, []byte(`{ invalid json`), 0o644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	_, err := resolveStartupOptions(startupOptions{
		configPath: configPath,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResolveStartupOptionsPreservesInitConfigFlag(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "config.json")

	resolved, err := resolveStartupOptions(startupOptions{
		city:       "Copenhagen",
		country:    "DK",
		configPath: configPath,
		initConfig: true,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !resolved.initConfig {
		t.Fatal("expected initConfig to be true")
	}

	if resolved.city != "Copenhagen" {
		t.Fatalf("expected city Copenhagen, got %q", resolved.city)
	}

	if resolved.country != "DK" {
		t.Fatalf("expected country DK, got %q", resolved.country)
	}
}

func TestResolveStartupOptionsPreservesFailFlag(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "config.json")

	resolved, err := resolveStartupOptions(startupOptions{
		city:       "Copenhagen",
		country:    "DK",
		configPath: configPath,
		fail:       true,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !resolved.fail {
		t.Fatal("expected fail to be true")
	}
}

func TestResolveStartupOptionsPreservesFakeFlag(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "config.json")

	resolved, err := resolveStartupOptions(startupOptions{
		city:       "Copenhagen",
		country:    "DK",
		configPath: configPath,
		fake:       true,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !resolved.fake {
		t.Fatal("expected fake to be true")
	}
}

func writeStartupTestConfig(t *testing.T, config meteoConfig.AppConfig) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.json")

	if err := meteoConfig.WriteConfig(path, config); err != nil {
		t.Fatalf("failed to write startup test config: %v", err)
	}

	return path
}
