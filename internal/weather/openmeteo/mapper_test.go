package openmeteo

import (
	"testing"

	"github.com/nfmdev/meteo/internal/domain"
)

func TestMapForecastResponseMapsDailyAndHourlyForecasts(t *testing.T) {
	t.Parallel()

	location := domain.Location{
		City:      "Copenhagen",
		Country:   "DK",
		Latitude:  55.6761,
		Longitude: 12.5683,
		Timezone:  "Europe/Copenhagen",
	}

	dto := forecastResponseDTO{
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

	report, err := mapForecastResponse(location, dto)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := len(report.Daily); got != 2 {
		t.Fatalf("expected 2 daily forecasts, got %d", got)
	}

	if got := len(report.Hourly); got != 2 {
		t.Fatalf("expected 2 hourly forecasts, got %d", got)
	}

	if got := report.Daily[0].Condition; got != "Partly Cloudy" {
		t.Fatalf("expected Partly Cloudy, got %q", got)
	}

	if got := report.Hourly[1].Condition; got != "Overcast" {
		t.Fatalf("expected Overcast, got %q", got)
	}
}
