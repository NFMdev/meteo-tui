package tui

import (
	"context"
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestModelCanBeCreatedWithLocationPreferencesLoader(t *testing.T) {
	t.Parallel()

	model := NewModelWithSearchAndPreferences(
		"Copenhagen",
		"DK",
		fakeWeatherLoader{},
		nil,
		fakeLocationPreferencesLoader{},
	)

	if model.locationPreferencesLoader == nil {
		t.Fatal("expected location preferences loader")
	}
}

var _ LocationPreferencesLoader = fakeLocationPreferencesLoader{}

type fakeLocationPreferencesLoader struct{}

func (fakeLocationPreferencesLoader) ListFavorites(
	ctx context.Context,
) ([]domain.SavedLocation, error) {
	return []domain.SavedLocation{}, nil
}

func (fakeLocationPreferencesLoader) AddFavorite(
	ctx context.Context,
	location domain.SavedLocation,
) error {
	return nil
}

func (fakeLocationPreferencesLoader) RemoveFavorite(
	ctx context.Context,
	location domain.SavedLocation,
) error {
	return nil
}

func (fakeLocationPreferencesLoader) SetDefaultLocation(
	ctx context.Context,
	location domain.SavedLocation,
) error {
	return nil
}

type fakeWeatherLoader struct{}

func (fakeWeatherLoader) GetWeather(
	ctx context.Context,
	city string,
	country string,
) (domain.WeatherReport, error) {
	return domain.WeatherReport{}, nil
}
