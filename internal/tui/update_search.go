package tui

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (m Model) updateSearchInputKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		return m.exitSearchMode(), nil

	case key.Matches(msg, m.keys.Submit):
		query := strings.TrimSpace(m.searchInput.Value())
		if query == "" {
			m.searchErr = errSearchQueryRequired
			return m, nil
		}

		m.mode = screenModeSearchResults
		m.searching = true
		m.searchErr = nil
		m.searchResults = nil
		m.selectedSearchResult = 0
		m.searchInput.Blur()

		return m, m.searchLocationsCmd(query)
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)

	return m, cmd
}

func (m Model) updateSearchResultKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		m.mode = screenModeSearchInput
		m.searching = false
		m.searchErr = nil
		return m, m.searchInput.Focus()

	case key.Matches(msg, m.keys.Up):
		if m.selectedSearchResult > 0 {
			m.selectedSearchResult--
		}
		return m, nil

	case key.Matches(msg, m.keys.Down):
		if m.selectedSearchResult < len(m.searchResults)-1 {
			m.selectedSearchResult++
		}
		return m, nil

	case key.Matches(msg, m.keys.Submit):
		result, ok := m.selectedLocationSearchResult()
		if !ok {
			m.searchErr = errSearchResultRequired
			return m, nil
		}

		m.city = result.Name
		m.country = result.CountryCode
		m.mode = screenModeDashboard
		m.loading = true
		m.err = nil
		m.selectedSearchResult = 0
		m.searchInput.SetValue("")
		m.searchInput.Blur()

		return m, m.loadWeatherCmd()
	}
	return m, m.loadWeatherCmd()
}
