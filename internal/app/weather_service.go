package app

import (
	"context"
	"time"

	"github.com/nfmdev/meteo/internal/domain"
)

type WeatherService interface {
	GetWeather(ctx context.Context, city string, country string) (domain.WeatherReport, error)
}

type FakeWeatherService struct{}

func NewFakeWeatherService() FakeWeatherService {
	return FakeWeatherService{}
}

func (s FakeWeatherService) GetWeather(
	ctx context.Context,
	city string,
	country string,
) (domain.WeatherReport, error) {
	now := time.Now()

	location := domain.Location{
		City:      city,
		Country:   country,
		Latitude:  57.048,
		Longitude: 9.9187,
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
			MaxWindKmh:      47.0,
		},
	}

	hourly := []domain.HourlyForecast{
		{
			Time:              beginningOfDay(now).Add(8 * time.Hour),
			TemperatureC:      15.8,
			FeelsLikeC:        11.2,
			Condition:         "Cloudy",
			WeatherCode:       3,
			HumidityPercent:   84,
			PrecipitationMM:   0.1,
			CloudCoverPercent: 92,
			WindSpeedKmh:      21.0,
		},
		{
			Time:              beginningOfDay(now).Add(9 * time.Hour),
			TemperatureC:      16.5,
			FeelsLikeC:        12.5,
			Condition:         "Cloudy",
			WeatherCode:       3,
			HumidityPercent:   82,
			PrecipitationMM:   0.1,
			CloudCoverPercent: 99,
			WindSpeedKmh:      22.5,
		},
		{
			Time:              beginningOfDay(now).Add(10 * time.Hour),
			TemperatureC:      17.8,
			FeelsLikeC:        14.2,
			Condition:         "Overcast",
			WeatherCode:       3,
			HumidityPercent:   80,
			PrecipitationMM:   0.2,
			CloudCoverPercent: 95,
			WindSpeedKmh:      24.0,
		},
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
