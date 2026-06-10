package openmeteo

import (
	"errors"
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestMapForecastResponseMapsCurrentWeather(t *testing.T) {
	location := testLocation()
	dto := testForecastResponseDTO()

	report, err := mapForecastResponse(location, dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.Location.City != "Copenhagen" {
		t.Fatalf("expected city Copenhagen, got %q", report.Location.City)
	}

	if report.Location.Country != "DK" {
		t.Fatalf("expected country DK, got %q", report.Location.Country)
	}

	if report.Location.Timezone != "Europe/Copenhagen" {
		t.Fatalf("expected timezone Europe/Copenhagen, got %q", report.Location.Timezone)
	}

	if report.Current.TemperatureC != 18.5 {
		t.Fatalf("expected current temperature 18.5, got %.1f", report.Current.TemperatureC)
	}

	if report.Current.FeelsLikeC != 17.9 {
		t.Fatalf("expected feels-like temperature 17.9, got %.1f", report.Current.FeelsLikeC)
	}

	if report.Current.Condition != "Partly Cloudy" {
		t.Fatalf("expected condition Partly Cloudy, got %q", report.Current.Condition)
	}

	if report.Current.WindSpeedKmh != 18.0 {
		t.Fatalf("expected wind speed 18.0, got %.1f", report.Current.WindSpeedKmh)
	}

	if report.Current.WindDirectionDeg != 240 {
		t.Fatalf("expected wind direction 240, got %d", report.Current.WindDirectionDeg)
	}
}

func TestMapForecastResponseMapsMetrics(t *testing.T) {
	location := testLocation()
	dto := testForecastResponseDTO()

	report, err := mapForecastResponse(location, dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.Metrics.HumidityPercent != 65 {
		t.Fatalf("expected humidity 65, got %d", report.Metrics.HumidityPercent)
	}

	if report.Metrics.PressureHPa != 1015.2 {
		t.Fatalf("expected pressure 1015.2, got %.1f", report.Metrics.PressureHPa)
	}

	if report.Metrics.PrecipitationMM != 0.0 {
		t.Fatalf("expected precipitation 0.0, got %.1f", report.Metrics.PrecipitationMM)
	}

	if report.Metrics.CloudCoverPercent != 40 {
		t.Fatalf("expected cloud cover 40, got %d", report.Metrics.CloudCoverPercent)
	}

	if report.Metrics.WindSpeedKmh != 18.0 {
		t.Fatalf("expected wind speed 18.0, got %.1f", report.Metrics.WindSpeedKmh)
	}
}

func TestMapForecastResponseMapsDailyForecasts(t *testing.T) {
	location := testLocation()
	dto := testForecastResponseDTO()

	report, err := mapForecastResponse(location, dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := len(report.Daily); got != 2 {
		t.Fatalf("expected 2 daily forecasts, got %d", got)
	}

	firstDay := report.Daily[0]

	if firstDay.Date.Year() != 2026 {
		t.Fatalf("expected year 2026, got %d", firstDay.Date.Year())
	}

	if int(firstDay.Date.Month()) != 6 {
		t.Fatalf("expected month 6, got %d", firstDay.Date.Month())
	}

	if firstDay.Date.Day() != 10 {
		t.Fatalf("expected day 10, got %d", firstDay.Date.Day())
	}

	if firstDay.MaxTemperatureC != 20.1 {
		t.Fatalf("expected max temperature 20.1, got %.1f", firstDay.MaxTemperatureC)
	}

	if firstDay.MinTemperatureC != 12.5 {
		t.Fatalf("expected min temperature 12.5, got %.1f", firstDay.MinTemperatureC)
	}

	if firstDay.Condition != "Partly Cloudy" {
		t.Fatalf("expected condition Partly Cloudy, got %q", firstDay.Condition)
	}

	if firstDay.PrecipitationMM != 0.0 {
		t.Fatalf("expected precipitation 0.0, got %.1f", firstDay.PrecipitationMM)
	}

	if firstDay.MaxWindKmh != 22.0 {
		t.Fatalf("expected max wind 22.0, got %.1f", firstDay.MaxWindKmh)
	}

	secondDay := report.Daily[1]

	if secondDay.Condition != "Rain" {
		t.Fatalf("expected second day condition Rain, got %q", secondDay.Condition)
	}

	if secondDay.PrecipitationMM != 2.4 {
		t.Fatalf("expected second day precipitation 2.4, got %.1f", secondDay.PrecipitationMM)
	}
}

func TestMapForecastResponseMapsHourlyForecasts(t *testing.T) {
	location := testLocation()
	dto := testForecastResponseDTO()

	report, err := mapForecastResponse(location, dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := len(report.Hourly); got != 2 {
		t.Fatalf("expected 2 hourly forecasts, got %d", got)
	}

	firstHour := report.Hourly[0]

	if firstHour.Time.Year() != 2026 {
		t.Fatalf("expected year 2026, got %d", firstHour.Time.Year())
	}

	if int(firstHour.Time.Month()) != 6 {
		t.Fatalf("expected month 6, got %d", firstHour.Time.Month())
	}

	if firstHour.Time.Day() != 10 {
		t.Fatalf("expected day 10, got %d", firstHour.Time.Day())
	}

	if firstHour.Time.Hour() != 8 {
		t.Fatalf("expected hour 8, got %d", firstHour.Time.Hour())
	}

	if firstHour.TemperatureC != 14.2 {
		t.Fatalf("expected temperature 14.2, got %.1f", firstHour.TemperatureC)
	}

	if firstHour.FeelsLikeC != 13.8 {
		t.Fatalf("expected feels-like 13.8, got %.1f", firstHour.FeelsLikeC)
	}

	if firstHour.Condition != "Partly Cloudy" {
		t.Fatalf("expected condition Partly Cloudy, got %q", firstHour.Condition)
	}

	if firstHour.HumidityPercent != 70 {
		t.Fatalf("expected humidity 70, got %d", firstHour.HumidityPercent)
	}

	if firstHour.CloudCoverPercent != 40 {
		t.Fatalf("expected cloud cover 40, got %d", firstHour.CloudCoverPercent)
	}

	secondHour := report.Hourly[1]

	if secondHour.Condition != "Overcast" {
		t.Fatalf("expected second hour condition Overcast, got %q", secondHour.Condition)
	}
}

func TestMapForecastResponseReturnsErrorForEmptyForecast(t *testing.T) {
	location := testLocation()

	dto := forecastResponseDTO{
		Timezone: "Europe/Copenhagen",
		Current: currentDTO{
			Time: "2026-06-10T12:00",
		},
		Daily:  dailyDTO{},
		Hourly: hourlyDTO{},
	}

	_, err := mapForecastResponse(location, dto)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrInvalidForecastResponse) {
		t.Fatalf("expected ErrInvalidForecastResponse, got %v", err)
	}
}

func TestParseOpenMeteoDate(t *testing.T) {
	parsed, err := parseOpenMeteoDate("Europe/Copenhagen", "2026-06-10")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if parsed.Year() != 2026 {
		t.Fatalf("expected year 2026, got %d", parsed.Year())
	}

	if int(parsed.Month()) != 6 {
		t.Fatalf("expected month 6, got %d", parsed.Month())
	}

	if parsed.Day() != 10 {
		t.Fatalf("expected day 10, got %d", parsed.Day())
	}

	if parsed.Hour() != 0 {
		t.Fatalf("expected hour 0, got %d", parsed.Hour())
	}
}

func TestParseOpenMeteoDateTime(t *testing.T) {
	parsed, err := parseOpenMeteoDateTime("Europe/Copenhagen", "2026-06-10T12:30")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if parsed.Year() != 2026 {
		t.Fatalf("expected year 2026, got %d", parsed.Year())
	}

	if int(parsed.Month()) != 6 {
		t.Fatalf("expected month 6, got %d", parsed.Month())
	}

	if parsed.Day() != 10 {
		t.Fatalf("expected day 10, got %d", parsed.Day())
	}

	if parsed.Hour() != 12 {
		t.Fatalf("expected hour 12, got %d", parsed.Hour())
	}

	if parsed.Minute() != 30 {
		t.Fatalf("expected minute 30, got %d", parsed.Minute())
	}
}

func TestWeatherCodeDescription(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		expected string
	}{
		{
			name:     "clear sky",
			code:     0,
			expected: "Clear sky",
		},
		{
			name:     "partly cloudy",
			code:     2,
			expected: "Partly Cloudy",
		},
		{
			name:     "overcast",
			code:     3,
			expected: "Overcast",
		},
		{
			name:     "rain",
			code:     61,
			expected: "Rain",
		},
		{
			name:     "thunderstorm",
			code:     95,
			expected: "Thunderstorm",
		},
		{
			name:     "unknown",
			code:     999,
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := weatherCodeDescription(tt.code)

			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func testLocation() domain.Location {
	return domain.Location{
		City:      "Copenhagen",
		Country:   "DK",
		Latitude:  55.6761,
		Longitude: 12.5683,
		Timezone:  "Europe/Copenhagen",
	}
}

func testForecastResponseDTO() forecastResponseDTO {
	return forecastResponseDTO{
		Timezone: "Europe/Copenhagen",
		Current: currentDTO{
			Time:                "2026-06-10T12:00",
			Temperature2m:       18.5,
			ApparentTemperature: 17.9,
			RelativeHumidity2m:  65,
			Precipitation:       0.0,
			WeatherCode:         2,
			CloudCover:          40,
			PressureMSL:         1015.2,
			WindSpeed10m:        18.0,
			WindDirection10m:    240,
		},
		Daily: dailyDTO{
			Time:             []string{"2026-06-10", "2026-06-11"},
			WeatherCode:      []int{2, 61},
			Temperature2mMax: []float64{20.1, 19.3},
			Temperature2mMin: []float64{12.5, 11.8},
			PrecipitationSum: []float64{0.0, 2.4},
			WindSpeed10mMax:  []float64{22.0, 28.0},
		},
		Hourly: hourlyDTO{
			Time:                []string{"2026-06-10T08:00", "2026-06-10T09:00"},
			Temperature2m:       []float64{14.2, 15.0},
			ApparentTemperature: []float64{13.8, 14.6},
			RelativeHumidity2m:  []int{70, 68},
			Precipitation:       []float64{0.0, 0.1},
			WeatherCode:         []int{2, 3},
			CloudCover:          []int{40, 60},
			PressureMSL:         []float64{1015.0, 1015.1},
			WindSpeed10m:        []float64{16.0, 17.5},
		},
	}
}
