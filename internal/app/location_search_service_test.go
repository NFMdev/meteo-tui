package app

import (
	"context"
	"errors"
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
)

var errTestLocationSearch = errors.New("location search failed")

func TestLocationSearchServiceReturnsResults(t *testing.T) {
	t.Parallel()

	expectedResults := []domain.LocationSearchResult{
		{
			Name:        "Copenhagen",
			Country:     "Denmark",
			CountryCode: "DK",
			Admin1:      "North Denmark",
			Latitude:    55.6761,
			Longitude:   12.5683,
			Timezone:    "Europe/Copenhagen",
		},
	}

	searcher := &fakeLocationSearcher{
		results: expectedResults,
	}

	service := NewLocationSearchService(searcher)

	results, err := service.SearchLocations(context.Background(), "Copenhagen")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !searcher.called {
		t.Fatal("expected searcher to be called")
	}
	if searcher.query != "Copenhagen" {
		t.Fatalf("expected query Copenhagen, got %q", searcher.query)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "Copenhagen" {
		t.Fatalf("expected Copenhagen, got %q", results[0].Name)
	}
	if results[0].CountryCode != "DK" {
		t.Fatalf("expected DK, got %q", results[0].CountryCode)
	}
}

func TestLocationSearchServiceTrimsQuery(t *testing.T) {
	t.Parallel()

	searcher := &fakeLocationSearcher{
		results: []domain.LocationSearchResult{},
	}

	service := NewLocationSearchService(searcher)

	_, err := service.SearchLocations(context.Background(), "  Copenhagen  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if searcher.query != "Copenhagen" {
		t.Fatalf("expected trimmed query Copenhagen, got %q", searcher.query)
	}
}

func TestLocationSearchServiceRejectsEmptyQuery(t *testing.T) {
	t.Parallel()

	searcher := &fakeLocationSearcher{}

	service := NewLocationSearchService(searcher)

	_, err := service.SearchLocations(context.Background(), "   ")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrLocationSearchQueryRequired) {
		t.Fatalf("expected ErrLocationSearchQueryRequired, got %v", err)
	}
	if searcher.called {
		t.Fatal("expected searcher not to be called")
	}
}

func TestLocationSearchServiceRequiresSearcher(t *testing.T) {
	t.Parallel()

	service := NewLocationSearchService(nil)

	_, err := service.SearchLocations(context.Background(), "Copenhagen")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrLocationSearcherRequired) {
		t.Fatalf("expected ErrLocationSearcherRequired, got %v", err)
	}
}

func TestLocationSearchServicePropagatesSearcherError(t *testing.T) {
	t.Parallel()

	searcher := &fakeLocationSearcher{
		err: errTestLocationSearch,
	}

	service := NewLocationSearchService(searcher)

	_, err := service.SearchLocations(context.Background(), "Copenhagen")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, errTestLocationSearch) {
		t.Fatalf("expected searcher error, got %v", err)
	}
}

func TestLocationSearchServiceAllowsEmptyResults(t *testing.T) {
	t.Parallel()

	searcher := &fakeLocationSearcher{
		results: []domain.LocationSearchResult{},
	}

	service := NewLocationSearchService(searcher)

	results, err := service.SearchLocations(context.Background(), "Unknown")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if results == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

var _ LocationSearcher = (*fakeLocationSearcher)(nil)

type fakeLocationSearcher struct {
	called  bool
	query   string
	results []domain.LocationSearchResult
	err     error
}

func (s *fakeLocationSearcher) Search(
	ctx context.Context,
	query string,
) ([]domain.LocationSearchResult, error) {
	s.called = true
	s.query = query
	if s.err != nil {
		return nil, s.err
	}
	return s.results, nil
}
