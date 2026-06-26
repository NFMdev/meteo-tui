package tui

import tea "charm.land/bubbletea/v2"

type screenMode int

const (
	screenModeDashboard screenMode = iota
	screenModeSearchInput
	screenModeSearchResults
)

func (m Model) enterSearchMode() (Model, tea.Cmd) {
	m.mode = screenModeSearchInput
	m.searching = false
	m.searchErr = nil
	m.searchResults = nil
	m.selectedSearchResult = 0
	m.searchInput.SetValue("")

	return m, m.searchInput.Focus()
}

func (m Model) exitSearchMode() Model {
	m.mode = screenModeDashboard
	m.searching = false
	m.searchErr = nil
	m.searchResults = nil
	m.selectedSearchResult = 0
	m.searchInput.SetValue("")
	m.searchInput.Blur()

	return m
}
