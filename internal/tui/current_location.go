package tui

import "github.com/nfmdev/meteo/internal/domain"

func (m Model) currentSavedLocation() (domain.SavedLocation, bool) {
	location := domain.SavedLocationFromLocation(m.report.Location)

	if domain.SavedLocationKey(location) == "" {
		return domain.SavedLocation{}, false
	}

	return location, true
}
