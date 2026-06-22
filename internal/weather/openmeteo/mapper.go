package openmeteo

import (
	"fmt"
	"time"

	"github.com/nfmdev/meteo/internal/domain"
)

func mapForecastResponse(
	location domain.Location,
	dto forecastResponseDTO,
) (domain.WeatherReport, error) {
	if len(dto.Daily.Time) == 0 || len(dto.Hourly.Time) == 0 {
		return domain.WeatherReport{}, ErrInvalidForecastResponse
	}

	if dto.Timezone != "" {
		location.Timezone = dto.Timezone
	}

	updatedAt, err := parseOpenMeteoDateTime(location.Timezone, dto.Current.Time)
	if err != nil {
		updatedAt = time.Now()
	}

	report := domain.WeatherReport{
		Location:  location,
		UpdatedAt: updatedAt,
		Source:    domain.NewFreshWeatherSource(domain.WeatherProviderOpenMeteo),
		Current: domain.CurrentWeather{
			TemperatureC:     dto.Current.Temperature2m,
			FeelsLikeC:       dto.Current.ApparentTemperature,
			Condition:        weatherCodeDescription(dto.Current.WeatherCode),
			WeatherCode:      dto.Current.WeatherCode,
			WindSpeedKmh:     dto.Current.WindSpeed10m,
			WindDirectionDeg: dto.Current.WindDirection10m,
		},
		Metrics: domain.WeatherMetrics{
			HumidityPercent:   dto.Current.RelativeHumidity2m,
			PressureHPa:       dto.Current.PressureMSL,
			PrecipitationMM:   dto.Current.Precipitation,
			CloudCoverPercent: dto.Current.CloudCover,
			WindSpeedKmh:      dto.Current.WindSpeed10m,
			WindDirectionDeg:  dto.Current.WindDirection10m,
		},
		Daily:  mapDailyForecasts(location.Timezone, dto.Daily),
		Hourly: mapHourlyForecasts(location.Timezone, dto.Hourly),
	}

	return report, nil
}

func mapDailyForecasts(timezone string, dto dailyDTO) []domain.DailyForecast {
	length := minInts(
		len(dto.Time),
		len(dto.WeatherCode),
		len(dto.Temperature2mMax),
		len(dto.Temperature2mMin),
		len(dto.PrecipitationSum),
		len(dto.WindSpeed10mMax),
	)

	forecasts := make([]domain.DailyForecast, 0, length)

	for i := range length {
		date, err := parseOpenMeteoDate(timezone, dto.Time[i])
		if err != nil {
			continue
		}

		code := dto.WeatherCode[i]

		forecasts = append(forecasts, domain.DailyForecast{
			Date:            date,
			MinTemperatureC: dto.Temperature2mMin[i],
			MaxTemperatureC: dto.Temperature2mMax[i],
			Condition:       weatherCodeDescription(code),
			WeatherCode:     code,
			PrecipitationMM: dto.PrecipitationSum[i],
			MaxWindKmh:      dto.WindSpeed10mMax[i],
		})
	}

	return forecasts
}

func mapHourlyForecasts(timezone string, dto hourlyDTO) []domain.HourlyForecast {
	length := minInts(
		len(dto.Time),
		len(dto.Temperature2m),
		len(dto.ApparentTemperature),
		len(dto.RelativeHumidity2m),
		len(dto.Precipitation),
		len(dto.WeatherCode),
		len(dto.CloudCover),
		len(dto.WindSpeed10m),
	)

	forecasts := make([]domain.HourlyForecast, 0, length)

	for i := range length {
		forecastTime, err := parseOpenMeteoDateTime(timezone, dto.Time[i])
		if err != nil {
			continue
		}

		code := dto.WeatherCode[i]

		forecasts = append(forecasts, domain.HourlyForecast{
			Time:              forecastTime,
			TemperatureC:      dto.Temperature2m[i],
			FeelsLikeC:        dto.ApparentTemperature[i],
			Condition:         weatherCodeDescription(code),
			WeatherCode:       code,
			HumidityPercent:   dto.RelativeHumidity2m[i],
			PrecipitationMM:   dto.Precipitation[i],
			CloudCoverPercent: dto.CloudCover[i],
			WindSpeedKmh:      dto.WindSpeed10m[i],
		})
	}

	return forecasts
}

func parseOpenMeteoDate(timezone string, value string) (time.Time, error) {
	location, err := loadTimeLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	return time.ParseInLocation("2006-01-02", value, location)
}

func parseOpenMeteoDateTime(timezone string, value string) (time.Time, error) {
	location, err := loadTimeLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	return time.ParseInLocation("2006-01-02T15:04", value, location)
}

func loadTimeLocation(timezone string) (*time.Location, error) {
	if timezone == "" || timezone == "auto" {
		return time.Local, nil
	}

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("load timezone %q: %w", timezone, err)
	}

	return location, nil
}

// There is no code for Windy weather
// Possible implementation when "Clear Sky" and High Wind Speed
func weatherCodeDescription(code int) string {
	switch code {
	case 0:
		return "Clear sky"
	case 1:
		return "Mainly Clear"
	case 2:
		return "Partly Cloudy"
	case 3:
		return "Overcast"
	case 45, 28:
		return "Fog"
	case 51, 53, 55:
		return "Drizzle"
	case 56, 57:
		return "Freezing Drizzle"
	case 61, 63, 65:
		return "Rain"
	case 66, 67:
		return "Frezing Rain"
	case 71, 73, 75:
		return "Snowfall"
	case 77:
		return "Snow Grains"
	case 80, 81, 82:
		return "Rain Showers"
	case 85, 86:
		return "Snow Showers"
	case 95:
		return "Thunderstorm"
	case 96, 99:
		return "Thunderstorm with hail"
	default:
		return "Unknown"
	}
}

func minInts(values ...int) int {
	if len(values) == 0 {
		return 0
	}

	min := values[0]
	for _, value := range values[1:] {
		if value < min {
			min = value
		}
	}

	return min
}
