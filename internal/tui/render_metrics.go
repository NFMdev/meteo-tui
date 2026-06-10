package tui

import (
	"fmt"
	"strings"
)

func (m Model) renderMetrics() string {
	metrics := m.report.Metrics

	content := strings.Join([]string{
		"Metrics",
		"",
		fmt.Sprintf("Humidity:\t%d%%", metrics.HumidityPercent),
		fmt.Sprintf("Pressure:\t%.1f hPa", metrics.PressureHPa),
		fmt.Sprintf("Pecipitation:\t%.1f mm", metrics.PrecipitationMM),
		fmt.Sprintf("Cloud cover:\t%d%%", metrics.CloudCoverPercent),
		fmt.Sprintf("Wind:\t\t%.1f Km/h", metrics.WindSpeedKmh),
	}, "\n")

	return panelStyle.Render(content)
}
