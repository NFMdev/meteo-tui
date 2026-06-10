package tui

import (
	"time"

	"github.com/nfmdev/meteo/internal/domain"
)

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
