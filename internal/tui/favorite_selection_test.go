package tui

import (
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestSelectedSavedLocationReturnsSelectedFavorite(t *testing.T) {
	t.Parallel()

	model := Model{
		favorites: []domain.SavedLocation{
			{
				Name:        "Copenhagen",
				CountryCode: "DK",
			},
			{
				Name:        "Madrid",
				CountryCode: "ES",
			},
		},
		selectedFavorite: 1,
	}

	location, ok := model.selectedSavedLocation()
	if !ok {
		t.Fatal("expected selected favorite")
	}

	if location.Name != "Madrid" {
		t.Fatalf("expected Madrid, got %q", location.Name)
	}

	if location.CountryCode != "ES" {
		t.Fatalf("expected ES, got %q", location.CountryCode)
	}
}

func TestSelectedSavedLocationReturnsFalseWhenNoFavorites(t *testing.T) {
	t.Parallel()

	model := Model{}

	_, ok := model.selectedSavedLocation()
	if ok {
		t.Fatal("expected no selected favorite")
	}
}

func TestSelectedSavedLocationReturnsFalseWhenIndexIsNegative(t *testing.T) {
	t.Parallel()

	model := Model{
		favorites: []domain.SavedLocation{
			{
				Name:        "Copenhagen",
				CountryCode: "DK",
			},
		},
		selectedFavorite: -1,
	}

	_, ok := model.selectedSavedLocation()
	if ok {
		t.Fatal("expected no selected favorite")
	}
}

func TestSelectedSavedLocationReturnsFalseWhenIndexIsOutOfBounds(t *testing.T) {
	t.Parallel()

	model := Model{
		favorites: []domain.SavedLocation{
			{
				Name:        "Copenhagen",
				CountryCode: "DK",
			},
		},
		selectedFavorite: 1,
	}

	_, ok := model.selectedSavedLocation()
	if ok {
		t.Fatal("expected no selected favorite")
	}
}
