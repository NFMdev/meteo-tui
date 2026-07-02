package tui

import "github.com/nfmdev/meteo/internal/domain"

func (m Model) selectedSearchResultSavedLocation() (domain.SavedLocation, bool) {
	result, ok := m.selectedLocationSearchResult()
	if !ok {
		return domain.SavedLocation{}, false
	}

	location := domain.SavedLocationFromSearchResult(result)

	if domain.SavedLocationKey(location) == "" {
		return domain.SavedLocation{}, false
	}

	return location, true
}
