package tui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (m Model) updateFavoritesKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		return m.exitFavoritesMode(), nil

	case key.Matches(msg, m.keys.Up):
		if m.selectedFavorite > 0 {
			m.selectedFavorite--
		}
		return m, nil

	case key.Matches(msg, m.keys.Down):
		if m.selectedFavorite < len(m.favorites)-1 {
			m.selectedFavorite++
		}
		return m, nil

	case key.Matches(msg, m.keys.Submit):
		location, ok := m.selectedSavedLocation()
		if !ok {
			m.favoritesErr = errFavoriteRequired
			return m, nil
		}

		m.city = location.Name
		m.country = location.CountryCode

		m.mode = screenModeDashboard
		m.loading = true
		m.err = nil
		m.statusMessage = ""
		m.selectedDay = 0

		m.favoritesLoading = false
		m.favoritesErr = nil
		m.selectedFavorite = 0

		return m, m.loadWeatherCmd()
	}

	return m, nil
}
