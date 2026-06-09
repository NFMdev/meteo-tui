package tui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/nfmdev/meteo/internal/domain"
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
	help := footerStyle.Render(m.help.View(m.keys))

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

func (m Model) renderHeader() string {
	location := m.report.Location

	title := titleStyle.Render(
		fmt.Sprintf("Meteo — %s, %s", location.City, location.Country),
	)

	updatedAt := subtitleStyle.Render(
		fmt.Sprintf("Updated: %s", m.report.UpdatedAt.Format("15:04:05")),
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

func (m Model) selectedDailyForecast() (domain.DailyForecast, bool) {
	if m.selectedDay < 0 || m.selectedDay >= len(m.report.Daily) {
		return domain.DailyForecast{}, false
	}

	return m.report.Daily[m.selectedDay], true
}

func (m Model) hourlyForecastForSelectedDay() []domain.HourlyForecast {
	selectedDay, ok := m.selectedDailyForecast()
	if !ok {
		return nil
	}

	hours := make([]domain.HourlyForecast, 0)

	for _, hour := range m.report.Hourly {
		if sameDay(hour.Time, selectedDay.Date) {
			hours = append(hours, hour)
		}
	}

	return hours
}

func sameDay(a time.Time, b time.Time) bool {
	aYear, aMonth, aDay := a.Date()
	bYear, bMonth, bDay := b.Date()

	return aYear == bYear && aMonth == bMonth && aDay == bDay
}
