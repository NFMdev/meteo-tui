package tui

import tea "charm.land/bubbletea/v2"

type screenMode int

const (
	screenModeDashboard screenMode = iota
	screenModeSearchInput
)

func (m Model) enterSearchMode() (Model, tea.Cmd) {
	m.mode = screenModeSearchInput
	m.searchErr = nil
	m.searchInput.SetValue("")
	m.searchInput.Focus()

	return m, m.searchInput.Focus()
}

func (m Model) exitSearchMode() Model {
	m.mode = screenModeDashboard
	m.searchErr = nil
	m.searchInput.SetValue("")
	m.searchInput.Blur()

	return m
}
