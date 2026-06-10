package app

import (
	"context"
	"errors"
	"time"

	"github.com/nfmdev/meteo/internal/domain"
)

var ErrWeatherUnavailable = errors.New("weather data is currently unavailable")

type WeatherService interface {
	GetWeather(ctx context.Context, city string, country string) (domain.WeatherReport, error)
}

type RealWeatherService struct {
	locationResolver LocationResolver
	forecastProvider ForecastProvider
}

// TODO: remove fake weather service before releasing v1.0
// Keeping it for development testing
type FakeWeatherService struct {
	locationResolver LocationResolver
	shouldFail       bool
}

func NewWeatherService(
	locationResolver LocationResolver,
	forecastProvider ForecastProvider,
) RealWeatherService {
	return RealWeatherService{
		locationResolver: locationResolver,
		forecastProvider: forecastProvider,
	}
}

func NewFakeWeatherService(locationResolver LocationResolver) FakeWeatherService {
	return FakeWeatherService{
		locationResolver: locationResolver,
	}
}

func NewFailingWeatherService(locationResolver LocationResolver) FakeWeatherService {
	return FakeWeatherService{
		locationResolver: locationResolver,
		shouldFail:       true,
	}
}

func (s RealWeatherService) GetWeather(
	ctx context.Context,
	city string,
	country string,
) (domain.WeatherReport, error) {
	location, err := s.locationResolver.Resolve(ctx, city, country)
	if err != nil {
		return domain.WeatherReport{}, err
	}

	report, err := s.forecastProvider.GetForecast(ctx, location)
	if err != nil {
		return domain.WeatherReport{}, err
	}

	return report, nil
}

func (s FakeWeatherService) GetWeather(
	ctx context.Context,
	city string,
	country string,
) (domain.WeatherReport, error) {
	select {
	case <-time.After(350 * time.Millisecond):
	case <-ctx.Done():
		return domain.WeatherReport{}, ctx.Err()
	}

	if s.shouldFail {
		return domain.WeatherReport{}, ErrWeatherUnavailable
	}

	now := time.Now()

	location, err := s.locationResolver.Resolve(ctx, city, country)
	if err != nil {
		return domain.WeatherReport{}, err
	}

	daily := []domain.DailyForecast{
		{
			Date:            beginningOfDay(now),
			MinTemperatureC: 12.3,
			MaxTemperatureC: 19.5,
			Condition:       "Cloudy",
			WeatherCode:     3,
			PrecipitationMM: 0.4,
			MaxWindKmh:      24.0,
		},
		{
			Date:            beginningOfDay(now.AddDate(0, 0, 1)),
			MinTemperatureC: 13.0,
			MaxTemperatureC: 20.1,
			Condition:       "Light rain",
			WeatherCode:     61,
			PrecipitationMM: 2.1,
			MaxWindKmh:      28.5,
		},
		{
			Date:            beginningOfDay(now.AddDate(0, 0, 2)),
			MinTemperatureC: 11.6,
			MaxTemperatureC: 17.4,
			Condition:       "Windy",
			WeatherCode:     3,
			PrecipitationMM: 0.0,
			MaxWindKmh:      39.0,
		},
	}

	hourly := []domain.HourlyForecast{
		fakeHour(now, 0, 8, 15.8, 12.2, "Cloudy", 84, 0.1, 92, 21.0),
		fakeHour(now, 0, 9, 16.3, 13.8, "Cloudy", 82, 0.1, 90, 22.5),
		fakeHour(now, 0, 10, 18.1, 15.4, "Overcast", 80, 0.2, 95, 24.0),

		fakeHour(now, 1, 8, 15.9, 12.3, "Light rain", 88, 0.6, 98, 25.0),
		fakeHour(now, 1, 9, 16.4, 12.8, "Light rain", 87, 0.7, 99, 27.0),
		fakeHour(now, 1, 10, 17.0, 13.3, "Rain", 89, 0.8, 100, 28.5),

		fakeHour(now, 2, 8, 13.8, 10.5, "Windy", 76, 0.0, 70, 32.0),
		fakeHour(now, 2, 9, 14.6, 10.4, "Windy", 74, 0.0, 65, 34.0),
		fakeHour(now, 2, 10, 15.5, 11.2, "Partly Cloudy", 72, 0.0, 55, 35.0),
	}

	return domain.WeatherReport{
		Location:  location,
		UpdatedAt: now,
		Current: domain.CurrentWeather{
			TemperatureC:     17.8,
			FeelsLikeC:       14.2,
			Condition:        "Overcast",
			WeatherCode:      3,
			WindSpeedKmh:     24.0,
			WindDirectionDeg: 240,
		},
		Metrics: domain.WeatherMetrics{
			HumidityPercent:   80,
			PressureHPa:       1012.4,
			PrecipitationMM:   0.2,
			CloudCoverPercent: 95,
			WindSpeedKmh:      24.0,
			WindDirectionDeg:  240,
		},
		Daily:  daily,
		Hourly: hourly,
	}, nil
}

func beginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()

	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func fakeHour(
	now time.Time,
	dayOffset int,
	hour int,
	temperatureC float64,
	feelsLikeC float64,
	condition string,
	humidityPercent int,
	precipitationMM float64,
	cloudCoverPercent int,
	windSpeedKmh float64,
) domain.HourlyForecast {
	return domain.HourlyForecast{
		Time:              beginningOfDay(now.AddDate(0, 0, dayOffset)).Add(time.Duration(hour) * time.Hour),
		TemperatureC:      temperatureC,
		FeelsLikeC:        feelsLikeC,
		Condition:         condition,
		WeatherCode:       3,
		HumidityPercent:   humidityPercent,
		PrecipitationMM:   precipitationMM,
		CloudCoverPercent: cloudCoverPercent,
		WindSpeedKmh:      windSpeedKmh,
	}
}
