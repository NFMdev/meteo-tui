package tui

import (
	"fmt"
	"strings"
)

func (m Model) renderHourlyForecast() string {
	selectedDay, ok := m.selectedDailyForecast()
	if !ok {
		return panelStyle.Render("Hourly Forecast\n\nNo selected day.")
	}

	hours := m.hourlyForecastForSelectedDay()

	lines := []string{
		fmt.Sprintf("Hourly Forecast — %s", selectedDay.Date.Format("Mon 02 Jun")),
		"",
	}

	if len(hours) == 0 {
		lines = append(lines, "No hourly forecast avaliable for this day.")
		return panelStyle.Render(strings.Join(lines, "\n"))
	}

	for _, hour := range hours {
		lines = append(
			lines,
			fmt.Sprintf("%s\t%.1f°C\tfeels %.1f°C\t%s\twind %.1f km/h",
				hour.Time.Format("15:04"),
				hour.TemperatureC,
				hour.FeelsLikeC,
				hour.Condition,
				hour.WindSpeedKmh,
			),
		)
	}

	return panelStyle.Render(strings.Join(lines, "\n"))
}
