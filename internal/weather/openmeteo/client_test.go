package openmeteo

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestClientGetForecastReturnsMappedWeatherReport(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		assertQueryValue(t, query.Get("latitude"), "55.676100", "latitude")
		assertQueryValue(t, query.Get("longitude"), "12.568300", "longitude")
		assertQueryValue(t, query.Get("forecast_days"), "7", "forecast_days")
		assertQueryValue(t, query.Get("timezone"), "Europe/Copenhagen", "timezone")

		if !strings.Contains(query.Get("current"), "temperature_2m") {
			t.Fatalf("expected current variables to include temperature_2m, got %q", query.Get("current"))
		}

		if !strings.Contains(query.Get("hourly"), "apparent_temperature") {
			t.Fatalf("expected hourly variables to include apparent_temperature, got %q", query.Get("hourly"))
		}

		if !strings.Contains(query.Get("daily"), "temperature_2m_max") {
			t.Fatalf("expected daily variables to include temperature_2m_max, got %q", query.Get("daily"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"latitude": 55.6761,
			"longitude": 12.5683,
			"generationtime_ms": 0.123,
			"utc_offset_seconds": 7200,
			"timezone": "Europe/Copenhagen",
			"timezone_abbreviation": "CEST",
			"current": {
				"time": "2026-06-10T12:00",
				"interval": 900,
				"temperature_2m": 18.5,
				"apparent_temperature": 17.9,
				"relative_humidity_2m": 65,
				"precipitation": 0.0,
				"weather_code": 2,
				"cloud_cover": 40,
				"pressure_msl": 1015.2,
				"wind_speed_10m": 18.0,
				"wind_direction_10m": 240
			},
			"daily": {
				"time": ["2026-06-10", "2026-06-11"],
				"weather_code": [2, 61],
				"temperature_2m_max": [20.1, 19.3],
				"temperature_2m_min": [12.5, 11.8],
				"precipitation_sum": [0.0, 2.4],
				"wind_speed_10m_max": [22.0, 28.0]
			},
			"hourly": {
				"time": ["2026-06-10T08:00", "2026-06-10T09:00"],
				"temperature_2m": [14.2, 15.0],
				"apparent_temperature": [13.8, 14.6],
				"relative_humidity_2m": [70, 68],
				"precipitation": [0.0, 0.1],
				"weather_code": [2, 3],
				"cloud_cover": [40, 60],
				"pressure_msl": [1015.0, 1015.1],
				"wind_speed_10m": [16.0, 17.5]
			}
		}`))
	}))
	defer server.Close()

	client := Client{
		client:  server.Client(),
		baseURL: server.URL,
	}

	report, err := client.GetForecast(context.Background(), testForecastLocation())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.Location.City != "Copenhagen" {
		t.Fatalf("expected city Copenhagen, got %q", report.Location.City)
	}

	if report.Current.Condition != "Partly Cloudy" {
		t.Fatalf("expected current condition Partly Cloudy, got %q", report.Current.Condition)
	}

	if got := len(report.Daily); got != 2 {
		t.Fatalf("expected 2 daily forecasts, got %d", got)
	}

	if got := len(report.Hourly); got != 2 {
		t.Fatalf("expected 2 hourly forecasts, got %d", got)
	}

	if report.Daily[1].Condition != "Rain" {
		t.Fatalf("expected second daily condition Rain, got %q", report.Daily[1].Condition)
	}

	if report.Hourly[1].Condition != "Overcast" {
		t.Fatalf("expected second hourly condition Overcast, got %q", report.Hourly[1].Condition)
	}
}

func TestClientGetForecastHandlesNon2xxResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "provider unavailable", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := Client{
		client:  server.Client(),
		baseURL: server.URL,
	}

	_, err := client.GetForecast(context.Background(), testForecastLocation())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrForecastUnavailable) {
		t.Fatalf("expected ErrForecastUnavailable, got %v", err)
	}
}

func TestClientGetForecastHandlesMalformedJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{ invalid json`))
	}))
	defer server.Close()

	client := Client{
		client:  server.Client(),
		baseURL: server.URL,
	}

	_, err := client.GetForecast(context.Background(), testForecastLocation())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "decode forecast response") {
		t.Fatalf("expected decode forecast response error, got %v", err)
	}
}

func TestClientGetForecastHandlesInvalidForecastResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"latitude": 55.6761,
			"longitude": 12.5683,
			"timezone": "Europe/Copenhagen",
			"current": {
				"time": "2026-06-10T12:00",
				"temperature_2m": 18.5
			},
			"daily": {},
			"hourly": {}
		}`))
	}))
	defer server.Close()

	client := Client{
		client:  server.Client(),
		baseURL: server.URL,
	}

	_, err := client.GetForecast(context.Background(), testForecastLocation())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrInvalidForecastResponse) {
		t.Fatalf("expected ErrInvalidForecastResponse, got %v", err)
	}
}

func assertQueryValue(t *testing.T, got string, expected string, name string) {
	t.Helper()

	if got != expected {
		t.Fatalf("expected %s=%q, got %q", name, expected, got)
	}
}

func testForecastLocation() domain.Location {
	return domain.Location{
		City:      "Copenhagen",
		Country:   "DK",
		Latitude:  55.6761,
		Longitude: 12.5683,
		Timezone:  "Europe/Copenhagen",
	}
}
