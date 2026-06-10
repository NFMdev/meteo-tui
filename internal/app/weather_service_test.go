package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nfmdev/meteo/internal/domain"
)

var (
	errTestLocation = errors.New("test location error")
	errTestForecast = errors.New("test forecast error")
)

func TestRealWeatherServiceGetWeatherResolvesLocationAndGetsForecast(t *testing.T) {
	t.Parallel()

	location := testAppLocation()
	expectedReport := testAppWeatherReport(location)

	locationResolver := &fakeLocationResolver{
		location: location,
	}

	forecastProvider := &fakeForecastProvider{
		report: expectedReport,
	}

	service := NewWeatherService(locationResolver, forecastProvider)

	report, err := service.GetWeather(context.Background(), "Copenhagen", "DK")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !locationResolver.called {
		t.Fatal("expected location resolver to be called")
	}

	if locationResolver.city != "Copenhagen" {
		t.Fatalf("expected resolver city Copenhagen, got %q", locationResolver.city)
	}

	if locationResolver.country != "DK" {
		t.Fatalf("expected resolver country DK, got %q", locationResolver.country)
	}

	if !forecastProvider.called {
		t.Fatal("expected forecast provider to be called")
	}

	if forecastProvider.location.City != "Copenhagen" {
		t.Fatalf("expected provider location city Copenhagen, got %q", forecastProvider.location.City)
	}

	if forecastProvider.location.Country != "DK" {
		t.Fatalf("expected provider location city DK, got %q", forecastProvider.location.Country)
	}

	if report.Location.City != expectedReport.Location.City {
		t.Fatalf("expected report city %q, got %q", expectedReport.Location.City, report.Location.City)
	}

	if report.Current.TemperatureC != expectedReport.Current.TemperatureC {
		t.Fatalf(
			"expected current temperature %.1f, got %.1f",
			expectedReport.Current.TemperatureC,
			report.Current.TemperatureC,
		)
	}
}

func TestRealWeatherServiceGetWeatherReturnsLocationResolverError(t *testing.T) {
	t.Parallel()

	locationResolver := &fakeLocationResolver{
		err: errTestLocation,
	}

	forecastProvider := &fakeForecastProvider{}

	service := NewWeatherService(locationResolver, forecastProvider)

	_, err := service.GetWeather(context.Background(), "Copenhagen", "DK")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, errTestLocation) {
		t.Fatalf("expected location error, got %v", err)
	}

	if !locationResolver.called {
		t.Fatal("expected location resolver to be called")
	}

	if forecastProvider.called {
		t.Fatal("expected forecast provider not to be called")
	}
}

func TestRealWeatherServiceGetWeatherReturnsForecastProviderError(t *testing.T) {
	t.Parallel()

	locationResolver := &fakeLocationResolver{
		location: testAppLocation(),
	}

	forecastProvider := &fakeForecastProvider{
		err: errTestForecast,
	}

	service := NewWeatherService(locationResolver, forecastProvider)

	_, err := service.GetWeather(context.Background(), "Copenhagen", "DK")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, errTestForecast) {
		t.Fatalf("expected forecast error, got %v", err)
	}

	if !locationResolver.called {
		t.Fatal("expected location resolver to be called")
	}

	if !forecastProvider.called {
		t.Fatal("expected forecast provider to be called")
	}
}

type fakeLocationResolver struct {
	location domain.Location
	err      error

	called  bool
	city    string
	country string
}

func (r *fakeLocationResolver) Resolve(
	ctx context.Context,
	city string,
	country string,
) (domain.Location, error) {
	r.called = true
	r.city = city
	r.country = country

	if r.err != nil {
		return domain.Location{}, r.err
	}

	return r.location, nil
}

type fakeForecastProvider struct {
	report domain.WeatherReport
	err    error

	called   bool
	location domain.Location
}

func (p *fakeForecastProvider) GetForecast(
	ctx context.Context,
	location domain.Location,
) (domain.WeatherReport, error) {
	p.called = true
	p.location = location

	if p.err != nil {
		return domain.WeatherReport{}, p.err
	}

	return p.report, nil
}

func testAppLocation() domain.Location {
	return domain.Location{
		City:      "Copenhagen",
		Country:   "DK",
		Latitude:  55.6761,
		Longitude: 12.5683,
		Timezone:  "Europe/Copenhagen",
	}
}

func testAppWeatherReport(location domain.Location) domain.WeatherReport {
	return domain.WeatherReport{
		Location:  location,
		UpdatedAt: time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
		Current: domain.CurrentWeather{
			TemperatureC:     18.5,
			FeelsLikeC:       17.9,
			Condition:        "Partly Cloudy",
			WeatherCode:      2,
			WindSpeedKmh:     18.0,
			WindDirectionDeg: 240,
		},
		Metrics: domain.WeatherMetrics{
			HumidityPercent:   65,
			PressureHPa:       1015.2,
			PrecipitationMM:   0.0,
			CloudCoverPercent: 40,
			WindSpeedKmh:      18.0,
			WindDirectionDeg:  240,
		},
		Daily: []domain.DailyForecast{
			{
				Date:            time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC),
				MinTemperatureC: 12.5,
				MaxTemperatureC: 20.1,
				Condition:       "Partly Cloudy",
				WeatherCode:     2,
				PrecipitationMM: 0.0,
				MaxWindKmh:      22.0,
			},
		},
		Hourly: []domain.HourlyForecast{
			{
				Time:              time.Date(2026, 6, 10, 8, 0, 0, 0, time.UTC),
				TemperatureC:      14.2,
				FeelsLikeC:        13.8,
				Condition:         "Partly Cloudy",
				WeatherCode:       2,
				HumidityPercent:   70,
				PrecipitationMM:   0.0,
				CloudCoverPercent: 40,
				WindSpeedKmh:      16.0,
			},
		},
	}
}
