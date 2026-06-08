package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m Model) View() tea.View {
	content := m.render()

	return tea.NewView(content)
}

func (m Model) render() string {
	header := titleStyle.Render(
		fmt.Sprintf("Meteo - %s, %s", m.city, m.country),
	)

	subtitle := subtitleStyle.Render(
		"v0.1 Block 1 - static TUI shell",
	)

	currentWeatherPanel := panelStyle.Render(strings.Join([]string{
		"Current Weather",
		"",
		"Real weather data will be added in the next blocks.",
		"",
		"Placeholder:",
		"Temperature: --",
		"Condition: --",
	}, "\n"))

	dashboardPanel := panelStyle.Render(strings.Join([]string{
		"Dashboard Areas",
		"",
		"• Header with location and last updated time",
		"• Current weather panel",
		"• Current metrics panel",
		"• Daily forecast list",
		"• Hourly forecast list for selected day",
	}, "\n"))

	footer := footerStyle.Render("q quit | ctrl+c quit")

	return appStyle.Render(strings.Join([]string{
		header,
		subtitle,
		"",
		currentWeatherPanel,
		"",
		dashboardPanel,
		"",
		footer,
	}, "\n"))
}
