package location

import (
	"net/url"
	"testing"
)

func TestBuildGeocodingURL(t *testing.T) {
	t.Parallel()

	rawURL, err := buildGeocodingURL(
		"https://geocoding-api.open-meteo.com/v1/search",
		"Copenhagen",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("expected valid URL, got %v", err)
	}

	query := parsedURL.Query()

	if got := query.Get("name"); got != "Copenhagen" {
		t.Fatalf("expected name Copenhagen, got %q", got)
	}

	if got := query.Get("count"); got != "10" {
		t.Fatalf("expected coun 10, got %q", got)
	}

	if got := query.Get("language"); got != "en" {
		t.Fatalf("expected language en, got %q", got)
	}

	if got := query.Get("format"); got != "json" {
		t.Fatalf("expected format json, got %q", got)
	}
}
