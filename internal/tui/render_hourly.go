package tui

import (
	"fmt"
)

func (m Model) renderHourlyForecast(width int, height int) string {
	selectedDay, ok := m.selectedDailyForecast()
	if !ok {
		return renderPanel("Hourly Forecast", []string{"No selected day."}, width, height)
	}

	hours := m.hourlyForecastForSelectedDay()

	lines := make([]string, 0, len(hours)+1)

	if len(hours) == 0 {
		lines = append(lines, "No hourly forecast avaliable for this day.")

		return renderPanel(
			fmt.Sprintf("Hourly Forecast — %s", selectedDay.Date.Format("Mon 02 Jun")),
			lines,
			width,
			height,
		)
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

	return renderPanel(
		fmt.Sprintf("Hourly Forecast — %s", selectedDay.Date.Format("Mon 02 Jun")),
		lines,
		width,
		height,
	)
}
