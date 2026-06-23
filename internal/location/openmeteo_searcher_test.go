package location

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestOpenMeteoSearcherSearchReturnsLocationCandidates(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertSearchQueryValue(t, r, "name", "copenhagen")
		assertSearchQueryValue(t, r, "count", "10")
		assertSearchQueryValue(t, r, "language", "en")
		assertSearchQueryValue(t, r, "format", "json")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"results": [
				{
					"id": 2624886,
					"name": "Copenhagen",
					"latitude": 55.6761,
					"longitude": 12.5683,
					"country": "Denmark",
					"country_code": "dk",
					"admin1": "Capital Region",
					"timezone": "Europe/Copenhagen",
					"population": 142937
				},
				{
					"id": 6543921,
					"name": "Copenhagen Portland",
					"latitude": 55.6835,
					"longitude": 12.5792,
					"country": "Denmark",
					"country_code": "DK",
					"admin1": "Capital Region",
					"timezone": "Europe/Copenhagen",
					"population": 0
				}
			]
		}`))
	}))
	defer server.Close()

	searcher := OpenMeteoSearcher{
		client:  server.Client(),
		baseURL: server.URL,
	}

	results, err := searcher.Search(t.Context(), "copenhagen")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	assertSearchResult(t, results[0], domain.LocationSearchResult{
		Name:        "Copenhagen",
		Country:     "Denmark",
		CountryCode: "DK",
		Admin1:      "Capital Region",
		Latitude:    55.6761,
		Longitude:   12.5683,
		Timezone:    "Europe/Copenhagen",
	})

	assertSearchResult(t, results[1], domain.LocationSearchResult{
		Name:        "Copenhagen Portland",
		Country:     "Denmark",
		CountryCode: "DK",
		Admin1:      "Capital Region",
		Latitude:    55.6835,
		Longitude:   12.5792,
		Timezone:    "Europe/Copenhagen",
	})
}

func TestOpenMeteoSearcherSearchTrimsQuery(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertSearchQueryValue(t, r, "name", "copenhagen")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results": []}`))
	}))
	defer server.Close()

	searcher := OpenMeteoSearcher{
		client:  server.Client(),
		baseURL: server.URL,
	}

	_, err := searcher.Search(t.Context(), "  copenhagen  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestOpenMeteoSearcherSearchRejectsEmptyQuery(t *testing.T) {
	t.Parallel()

	searcher := NewOpenMeteoSearcher(nil)

	_, err := searcher.Search(t.Context(), "   ")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrSearchQueryRequired) {
		t.Fatalf("expected ErrSearchQueryRequired, got %v", err)
	}
}

func TestOpenMeteoSearcherSearchReturnsEmptySliceWhenNoResults(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results": []}`))
	}))
	defer server.Close()

	searcher := OpenMeteoSearcher{
		client:  server.Client(),
		baseURL: server.URL,
	}

	results, err := searcher.Search(t.Context(), "unknown")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("expected empty results, got %d", len(results))
	}
}

func TestOpenMeteoSearcherSearchHandlesMissingResultsField(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	searcher := OpenMeteoSearcher{
		client:  server.Client(),
		baseURL: server.URL,
	}

	results, err := searcher.Search(t.Context(), "unknown")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(results) != 0 {
		t.Fatalf("expected empty results, got %d", len(results))
	}
}

func TestOpenMeteoSearcherSearchSkipsInvalidResults(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"results": [
				{
					"name": "",
					"country_code": "DK"
				},
				{
					"name": "Copenhagen",
					"country_code": ""
				},
				{
					"name": "Copenhagen",
					"country": "Denmark",
					"country_code": "DK",
					"admin1": "Capital Region",
					"latitude": 55.6761,
					"longitude": 12.5683,
					"timezone": "Europe/Copenhagen"
				}
			]
		}`))
	}))
	defer server.Close()

	searcher := OpenMeteoSearcher{
		client:  server.Client(),
		baseURL: server.URL,
	}

	results, err := searcher.Search(t.Context(), "copenhagen")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 valid result, got %d", len(results))
	}

	if results[0].Name != "Copenhagen" {
		t.Fatalf("expected Copenhagen, got %q", results[0].Name)
	}
}

func TestOpenMeteoSearcherSearchHandlesNon2xxResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "upstream unavailable", http.StatusBadGateway)
	}))
	defer server.Close()

	searcher := OpenMeteoSearcher{
		client:  server.Client(),
		baseURL: server.URL,
	}

	_, err := searcher.Search(t.Context(), "copenhagen")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOpenMeteoSearcherSearchHandlesMalformedJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{ invalid json`))
	}))
	defer server.Close()

	searcher := OpenMeteoSearcher{
		client:  server.Client(),
		baseURL: server.URL,
	}

	_, err := searcher.Search(t.Context(), "copenhagen")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func assertSearchQueryValue(
	t *testing.T,
	request *http.Request,
	key string,
	expected string,
) {
	t.Helper()

	got := request.URL.Query().Get(key)
	if got != expected {
		t.Fatalf("expected query param %s=%q, got %q", key, expected, got)
	}
}

func assertSearchResult(
	t *testing.T,
	got domain.LocationSearchResult,
	expected domain.LocationSearchResult,
) {
	t.Helper()

	if got.Name != expected.Name {
		t.Fatalf("expected name %q, got %q", expected.Name, got.Name)
	}

	if got.Country != expected.Country {
		t.Fatalf("expected country %q, got %q", expected.Country, got.Country)
	}

	if got.CountryCode != expected.CountryCode {
		t.Fatalf("expected country code %q, got %q", expected.CountryCode, got.CountryCode)
	}

	if got.Admin1 != expected.Admin1 {
		t.Fatalf("expected admin1 %q, got %q", expected.Admin1, got.Admin1)
	}

	if got.Latitude != expected.Latitude {
		t.Fatalf("expected latitude %.4f, got %.4f", expected.Latitude, got.Latitude)
	}

	if got.Longitude != expected.Longitude {
		t.Fatalf("expected longitude %.4f, got %.4f", expected.Longitude, got.Longitude)
	}

	if got.Timezone != expected.Timezone {
		t.Fatalf("expected timezone %q, got %q", expected.Timezone, got.Timezone)
	}
}
