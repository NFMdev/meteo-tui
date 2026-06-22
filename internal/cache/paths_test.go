package cache

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestDefaultForecastCacheDirUsesXDGCacheHome(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "/tmp/meteo-cache-test")

	got, err := DefaultForecastCacheDir()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := filepath.Join("/tmp/meteo-cache-test", AppName, ForecastCacheDir)
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestResolveForecastCacheDirUsesDefaultWhenCustomDirIsEmpty(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", "/tmp/meteo-cache-test")

	got, err := ResolveForecastCacheDir("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := filepath.Join("/tmp/meteo-cache-test", AppName, ForecastCacheDir)
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestResolveForecastCacheDirTrimsAndCleansCustomDir(t *testing.T) {
	t.Parallel()

	got, err := ResolveForecastCacheDir("  /tmp/meteo-cache/../meteo-cache  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := filepath.Clean("/tmp/meteo-cache/../meteo-cache")
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestResolveForecastCacheDirExpandsHomeDir(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	got, err := ResolveForecastCacheDir("~/meteo-cache")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := filepath.Join(homeDir, "meteo-cache")
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestCacheFilePath(t *testing.T) {
	t.Parallel()

	got, err := CacheFilePath("/tmp/meteo-cache-test", "copenhagen_dk")
	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	expected := filepath.Join("/tmp/meteo-cache-test", "copenhagen_dk.json")
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestForecastCacheFilePath(t *testing.T) {
	t.Parallel()

	got, err := ForecastCacheFilePath("/tmp/meteo-cache-test", "Copenhagen", "DK")
	if err != nil {
		t.Fatalf("expectedo no errors, got %v", err)
	}

	expected := filepath.Join("/tmp/meteo-cache-test", "copenhagen_dk.json")
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestCacheFilePathRejectsEmptyCacheDir(t *testing.T) {
	t.Parallel()

	_, err := CacheFilePath("", "copenhagen_dk")

	if !errors.Is(err, ErrCacheDirRequired) {
		t.Fatalf("expected ErrCacheDirRequired, got %v", err)
	}
}

func TestCacheFilePathRejectsEmptyKey(t *testing.T) {
	t.Parallel()

	_, err := CacheFilePath("/tmp/meteo-cache-test", "")

	if !errors.Is(err, ErrCacheKeyRequired) {
		t.Fatalf("expected ErrCacheKeyRequired, got %v", err)
	}
}
