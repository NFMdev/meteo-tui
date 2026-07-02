package tui

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (m Model) updateDashboardKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Search):
		return m.enterSearchMode()

	case key.Matches(msg, m.keys.Favorites):
		return m.enterFavoritesMode()

	case key.Matches(msg, m.keys.AddFavorite):
		location, ok := m.currentSavedLocation()
		if !ok {
			m.statusMessage = ""
			m.err = errCurrentLocationRequired
			return m, nil
		}

		m.statusMessage = ""
		return m, m.addFavoriteCmd(location)

	case key.Matches(msg, m.keys.SetDefault):
		location, ok := m.currentSavedLocation()
		if !ok {
			m.statusMessage = ""
			return m, m.addFavoriteCmd(location)
		}

		m.statusMessage = ""
		return m, m.setDefaultLocationCmd(location)

	case key.Matches(msg, m.keys.Up):
		m.selectPreviousDay()
		m.rebuildViewportContent()

	case key.Matches(msg, m.keys.Down):
		m.selectNextDay()
		m.rebuildViewportContent()

	case key.Matches(msg, m.keys.ScrollUp):
		if m.layoutMode() == layoutModeCompactScrollable {
			m.viewport.ScrollUp(3)
		}

	case key.Matches(msg, m.keys.ScrollDown):
		if m.layoutMode() == layoutModeCompactScrollable {
			m.viewport.ScrollDown(3)
		}

	case key.Matches(msg, m.keys.ScrollTop):
		if m.layoutMode() == layoutModeCompactScrollable {
			m.viewport.GotoTop()
		}

	case key.Matches(msg, m.keys.ScrollBottom):
		if m.layoutMode() == layoutModeCompactScrollable {
			m.viewport.GotoBottom()
		}

	case key.Matches(msg, m.keys.Help):
		m.showHelp = !m.showHelp
		m.help.ShowAll = m.showHelp

	case key.Matches(msg, m.keys.Refresh):
		m.loading = true
		m.err = nil
		m.selectedDay = 0

		return m, tea.Batch(
			m.spinner.Tick,
			m.loadWeatherCmd(),
		)
	}
	return m, nil
}

func (m *Model) selectPreviousDay() {
	if m.loading || m.err != nil {
		return
	}

	if m.selectedDay > 0 {
		m.selectedDay--
	}
}

func (m *Model) selectNextDay() {
	if m.loading || m.err != nil {
		return
	}

	if m.selectedDay < len(m.report.Daily)-1 {
		m.selectedDay++
	}
}
