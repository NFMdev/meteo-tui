package tui

import "strings"

func (m *Model) configureViewport() {
	m.viewport.SetWidth(m.contentWidth())
	m.viewport.SetHeight(m.compactViewportHeight())
}

func (m *Model) rebuildViewportContent() {
	if m.loading || m.err != nil {
		return
	}

	m.configureViewport()

	content := m.compactScrollableContent()

	m.viewport.SetContent(content)
}

func (m Model) compactScrollableContent() string {
	return strings.Join([]string{
		m.renderCurrentWeather(m.panelWidth(), 7),
		"",
		m.renderMetrics(m.panelWidth(), 9),
		"",
		m.renderDailyForecast(m.panelWidth(), 10),
		"",
		m.renderHourlyForecast(m.panelWidth(), 16),
	}, "\n")
}
