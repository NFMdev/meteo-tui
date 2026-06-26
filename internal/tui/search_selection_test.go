package tui

import (
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestSelectedLocationSearchResultReturnsSelectedResult(t *testing.T) {
	t.Parallel()

	model := Model{
		searchResults: []domain.LocationSearchResult{
			{
				Name:        "Copenhagen",
				CountryCode: "DK",
			},
			{
				Name:        "Madrid",
				CountryCode: "ES",
			},
		},
		selectedSearchResult: 1,
	}

	result, ok := model.selectedLocationSearchResult()
	if !ok {
		t.Fatal("expected selected result")
	}

	if result.Name != "Madrid" {
		t.Fatalf("expected Madrid, got %q", result.Name)
	}

	if result.CountryCode != "ES" {
		t.Fatalf("expected ES, got %q", result.CountryCode)
	}
}

func TestSelectedLocationSearchResultReturnsFalseWhenNoResults(t *testing.T) {
	t.Parallel()

	model := Model{}

	_, ok := model.selectedLocationSearchResult()
	if ok {
		t.Fatal("expected no selected result")
	}
}

func TestSelectedLocationSearchResultReturnsFalseWhenIndexIsNegative(t *testing.T) {
	t.Parallel()

	model := Model{
		searchResults: []domain.LocationSearchResult{
			{
				Name:        "Copenhagen",
				CountryCode: "DK",
			},
		},
		selectedSearchResult: -1,
	}

	_, ok := model.selectedLocationSearchResult()
	if ok {
		t.Fatal("expected no selected result")
	}
}

func TestSelectedLocationSearchResultReturnsFalseWhenIndexIsOutOfBounds(t *testing.T) {
	t.Parallel()

	model := Model{
		searchResults: []domain.LocationSearchResult{
			{
				Name:        "Copenhagen",
				CountryCode: "DK",
			},
		},
		selectedSearchResult: 1,
	}

	_, ok := model.selectedLocationSearchResult()
	if ok {
		t.Fatal("expected no selected result")
	}
}
