package tui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m Model) View() tea.View {
	return tea.NewView(m.render())
}

func (m Model) render() string {
	header := m.renderHeader()
	current := m.renderCurrentWeather()
	metrics := m.renderMetrics()
	daily := m.renderDailyForecast()
	hourly := m.renderHourlyForecast()
	footer := footerStyle.Render("q quit | ctrl+c quit")

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
		footer,
	}, "\n"))
}

func (m Model) renderHeader() string {
	location := m.report.Location

	title := titleStyle.Render(
		fmt.Sprintf("Meteo — %s, %s", location.City, location.Country),
	)

	updatedAt := subtitleStyle.Render(
		fmt.Sprintf("Updated: %s", m.report.UpdatedAt.Format("15:04")),
	)

	return strings.Join([]string{
		title,
		updatedAt,
	}, "\n")
}

func (m Model) renderCurrentWeather() string {
	current := m.report.Current

	content := strings.Join([]string{
		"Current Weather",
		"",
		fmt.Sprintf("%.1f°C\t\t%s", current.TemperatureC, current.Condition),
		fmt.Sprintf("Feels Like\t%.1f°C", current.FeelsLikeC),
		fmt.Sprintf("Wind %.1f Km/h from %d°", current.WindSpeedKmh, current.WindDirectionDeg),
	}, "\n")

	return panelStyle.Render(content)
}

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

func (m Model) renderDailyForecast() string {
	lines := []string{
		"Daily Forecast",
		"",
	}

	for _, day := range m.report.Daily {
		lines = append(
			lines,
			fmt.Sprintf(
				"%s\t%.1f°C / %.1f°C\t%s\train %.1f mm",
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

func (m Model) renderHourlyForecast() string {
	lines := []string{
		"Hourly Forecast",
		"",
	}

	for _, hour := range m.report.Hourly {
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
