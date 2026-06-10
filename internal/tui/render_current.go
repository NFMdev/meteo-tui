package tui

import "fmt"

func (m Model) renderCurrentWeather() string {
	current := m.report.Current

	lines := []string{
		"Current Weather",
		"",
		fmt.Sprintf("%.1f°C\t\t%s", current.TemperatureC, current.Condition),
		fmt.Sprintf("Feels Like\t%.1f°C", current.FeelsLikeC),
		fmt.Sprintf("Wind %.1f Km/h from %d°", current.WindSpeedKmh, current.WindDirectionDeg),
	}

	return panelStyle.
		Width(m.panelWidth()).
		Render(joinTruncatedLines(lines, m.innerPanelWidth()))
}
