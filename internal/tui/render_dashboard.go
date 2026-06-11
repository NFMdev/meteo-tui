package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
)

func (m Model) renderDashboardGrid() string {
	layout := m.dashboardLayout()

	header := m.renderHeader()

	leftColumn := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderCurrentWeather(layout.leftWidth, layout.currentHeight),
		m.renderMetrics(layout.leftWidth, layout.metricsHeight),
		m.renderDailyForecast(layout.leftWidth, layout.dailyHeight),
	)

	rightColumn := m.renderHourlyForecast(layout.rightWidth, layout.hourlyHeight)

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftColumn,
		blank(layout.gap),
		rightColumn,
	)

	help := footerStyle.Render(m.help.View(m.keys))

	return appStyle.Render(strings.Join([]string{
		header,
		"",
		body,
		"",
		help,
	}, "\n"))
}

func (m Model) renderCompactScrollableDashboard() string {
	header := m.renderHeader()

	help := footerStyle.Render(
		m.help.View(m.keys),
	)

	return appStyle.Render(strings.Join([]string{
		header,
		"",
		m.viewport.View(),
		"",
		help,
	}, "\n"))
}
