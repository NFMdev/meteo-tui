package openmeteo

import (
	"net/url"
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestBuildForecastURL(t *testing.T) {
	t.Parallel()

	location := domain.Location{
		City:      "Copenhagen",
		Country:   "DK",
		Latitude:  55.6761,
		Longitude: 12.5683,
		Timezone:  "Europe/Copenhagen",
	}

	rawURL, err := buildForecastURL(
		"https://api.open-meteo.com/v1/forecast",
		location,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("expected valid URL, got %v", err)
	}

	query := parsedURL.Query()

	assertQueryValue(t, query, "latitude", "55.676100")
	assertQueryValue(t, query, "longitude", "12.568300")
	assertQueryValue(t, query, "forecast_days", "7")
	assertQueryValue(t, query, "timezone", "Europe/Copenhagen")

	if got := query.Get("current"); got == "" {
		t.Fatal("expected current variables")
	}

	if got := query.Get("hourly"); got == "" {
		t.Fatal("expected hourly variables")
	}

	if got := query.Get("daily"); got == "" {
		t.Fatal("expected daily variables")
	}
}

func assertQueryValue(
	t *testing.T,
	query url.Values,
	key string,
	expected string,
) {
	t.Helper()

	if got := query.Get(key); got != expected {
		t.Fatalf("expected %s=%q, got %q", key, expected, got)
	}
}
