package cache

import (
	"errors"
	"testing"
)

func TestForecastCacheKeyGeneratesStableKey(t *testing.T) {
	t.Parallel()

	key, err := ForecastCacheKey("Copenhagen", "DK")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if key != "copenhagen_dk" {
		t.Fatalf("expected copenhagen_dk, got %q", key)
	}
}

func TestForecastCacheKeyNormalizesInput(t *testing.T) {
	t.Parallel()

	key, err := ForecastCacheKey("   COpenHagEn  ", "dK   ")
	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	if key != "copenhagen_dk" {
		t.Fatalf("expected copenhagen_dk, got %q", key)
	}
}

func TestForecastCacheKeyHandlesSpaces(t *testing.T) {
	t.Parallel()

	key, err := ForecastCacheKey("New York", "US")
	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	if key != "new_york_us" {
		t.Fatalf("expected new_york_us, got %q", key)
	}
}

func TestForecastCacheKeyHandlesUnicodeLetters(t *testing.T) {
	t.Parallel()

	key, err := ForecastCacheKey("São Paulo", "BR")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if key != "são_paulo_br" {
		t.Fatalf("expected são_paulo_br, got %q", key)
	}
}

func TestForecastCacheKeyCollapsesSeparators(t *testing.T) {
	t.Parallel()

	key, err := ForecastCacheKey("New---York", "US")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if key != "new_york_us" {
		t.Fatalf("expected new_york_us, got %q", key)
	}
}

func TestForecastCacheKeyRejectsEmptyCity(t *testing.T) {
	t.Parallel()

	_, err := ForecastCacheKey("", "DK")
	if !errors.Is(err, ErrCacheCityRequired) {
		t.Fatalf("expected ErrCacheCityRequired, got %v", err)
	}
}

func TestForecastCacheKeyRejectsEmptyCountry(t *testing.T) {
	t.Parallel()

	_, err := ForecastCacheKey("Copenhagen", "")
	if !errors.Is(err, ErrCacheCountryRequired) {
		t.Fatalf("expected ErrCacheCountryRequired, got %v", err)
	}
}

func TestForecastCacheKeyRejectsInvalidCountry(t *testing.T) {
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
			country: "DKK",
		},
		{
			name:    "number",
			country: "D1",
		},
		{
			name:    "symbol",
			country: "D-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ForecastCacheKey("Copenhagen", tt.country)

			if !errors.Is(err, ErrInvalidCacheCountry) {
				t.Fatalf("expected ErrInvalidCacheCountry, got %v", err)
			}
		})
	}
}

func TestForecastCacheKeyRejectsInvalidKey(t *testing.T) {
	t.Parallel()

	_, err := ForecastCacheKey("---", "DK")

	if !errors.Is(err, ErrInvalidCacheKey) {
		t.Fatalf("expected ErrInvalidcacheKey, got %v", err)
	}
}
