package tui

import "github.com/nfmdev/meteo/internal/domain"

func (m Model) selectedSavedLocation() (domain.SavedLocation, bool) {
	if len(m.favorites) == 0 {
		return domain.SavedLocation{}, false
	}

	if m.selectedFavorite < 0 || m.selectedFavorite >= len(m.favorites) {
		return domain.SavedLocation{}, false
	}

	return m.favorites[m.selectedFavorite], true
}
