package tui

import "strings"

func (m Model) renderDashboard() string {
	header := m.renderHeader()
	current := m.renderCurrentWeather()
	metrics := m.renderMetrics()
	daily := m.renderDailyForecast()
	hourly := m.renderHourlyForecast()
	help := footerStyle.Render(truncateText(m.help.View(m.keys), m.contentWidth()))

	return appStyle.Render(strings.Join([]string{
		header,
		"",
		current,
		"",
		metrics,
		"",
		daily,
		"",
		hourly,
		"",
		help,
	}, "\n"))
}
