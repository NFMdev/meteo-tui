package tui

import (
	"fmt"
	"strings"
)

func (m Model) renderDailyForecast() string {
	lines := []string{
		"Daily Forecast",
		"",
	}

	for index, day := range m.report.Daily {
		cursor := " "
		if index == m.selectedDay {
			cursor = ">"
		}

		lines = append(
			lines,
			fmt.Sprintf(
				"%s\t%s\t%.1f°C / %.1f°C\t%s\train %.1f mm",
				cursor,
				day.Date.Format("Mon 02 Jun"),
				day.MaxTemperatureC,
				day.MinTemperatureC,
				day.Condition,
				day.PrecipitationMM,
			),
		)
	}

	return panelStyle.Render(strings.Join(lines, "\n"))
}
