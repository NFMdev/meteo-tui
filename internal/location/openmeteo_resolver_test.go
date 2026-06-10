package location

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenMeteoResolverResolveReturnsMatchingLocation(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		if got := query.Get("name"); got != "Copenhagen" {
			t.Fatalf("expected name=Copenhagen, got %q", got)
		}

		if got := query.Get("count"); got != "10" {
			t.Fatalf("expected count=10, got %q", got)
		}

		if got := query.Get("language"); got != "en" {
			t.Fatalf("expected language=en, got %q", got)
		}

		if got := query.Get("format"); got != "json" {
			t.Fatalf("expected format=json, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"results": [
				{
					"id": 123,
					"name": "Copenhagen",
					"latitude": 55.6761,
					"longitude": 12.5683,
					"country_code": "DK",
					"timezone": "Europe/Copenhagen",
					"country": "Denmark",
					"admin1": "North Denmark"
				}
			]
		}`))
	}))
	defer server.Close()

	resolver := OpenMeteoResolver{
		client:  server.Client(),
		baseURL: server.URL,
	}

	location, err := resolver.Resolve(context.Background(), "Copenhagen", "DK")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if location.City != "Copenhagen" {
		t.Fatalf("expected city Copenhagen, got %q", location.City)
	}

	if location.Country != "DK" {
		t.Fatalf("expected country DK, got %q", location.Country)
	}

	if location.Latitude != 55.6761 {
		t.Fatalf("expected latitude 55.6761, got %f", location.Latitude)
	}

	if location.Longitude != 12.5683 {
		t.Fatalf("expected longitude 12.5683, got %f", location.Longitude)
	}

	if location.Timezone != "Europe/Copenhagen" {
		t.Fatalf("expected timezone Europe/Copenhagen, got %q", location.Timezone)
	}
}

func TestOpenMeteoResolverResolveFiltersByCountry(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"results": [
				{
					"name": "Copenhagen",
					"latitude": 55.6761,
					"longitude": 12.5683,
					"country_code": "DK",
					"timezone": "Europe/Copenhagen"
				},
				{
					"name": "Copenhagen",
					"latitude": 10.000,
					"longitude": 20.000,
					"country_code": "XX",
					"timezone": "Etc/UTC"
				}
			]
		}`))
	}))
	defer server.Close()

	resolver := OpenMeteoResolver{
		client:  server.Client(),
		baseURL: server.URL,
	}

	location, err := resolver.Resolve(context.Background(), "Copenhagen", "XX")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if location.Country != "XX" {
		t.Fatalf("expected country XX, got %q", location.Country)
	}

	if location.Timezone != "Etc/UTC" {
		t.Fatalf("expected timezone Etc/UTC, got %q", location.Timezone)
	}
}

func TestOpenMeteoResolverResolveReturnsLocationNotFound(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"results": [
				{
					"name": "Copenhagen",
					"latitude": 55.6761,
					"longitude": 12.5683,
					"country_code": "DK",
					"timezone": "Europe/Copenhagen"
				}
			]
		}`))
	}))
	defer server.Close()

	resolver := OpenMeteoResolver{
		client:  server.Client(),
		baseURL: server.URL,
	}

	_, err := resolver.Resolve(context.Background(), "Copenhagen", "ES")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrLocationNotFound) {
		t.Fatalf("expected ErrLocationNotFound, got %v", err)
	}
}

func TestOpenMeteoResolverResolveHandlesNon2xxResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "provider unavailable", http.StatusInternalServerError)
	}))
	defer server.Close()

	resolver := OpenMeteoResolver{
		client:  server.Client(),
		baseURL: server.URL,
	}

	_, err := resolver.Resolve(context.Background(), "Copenhagen", "DK")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "stats 500") {
		t.Fatalf("expected stats 500 error, got %v", err)
	}
}

func TestOpenMeteoResolverResolveHandlesMalformedJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{ invalid json`))
	}))
	defer server.Close()

	resolver := OpenMeteoResolver{
		client:  server.Client(),
		baseURL: server.URL,
	}

	_, err := resolver.Resolve(context.Background(), "Copenhagen", "DK")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "decode geocoding response") {
		t.Fatalf("expected decode error, got %v", err)
	}
}

func TestOpenMeteoResolverResolveValidatesInput(t *testing.T) {
	t.Parallel()

	resolver := OpenMeteoResolver{
		client:  http.DefaultClient,
		baseURL: "http://example.com",
	}

	_, err := resolver.Resolve(context.Background(), "", "DK")
	if !errors.Is(err, ErrCityRequired) {
		t.Fatalf("expected ErrCityRequired, got %v", err)
	}

	_, err = resolver.Resolve(context.Background(), "Copenhagen", "")
	if !errors.Is(err, ErrCountryRequired) {
		t.Fatalf("expected ErrCountryRequired, got %v", err)
	}
}
