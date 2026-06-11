package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func writeTestConfigFile(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.json")

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	return path
}
