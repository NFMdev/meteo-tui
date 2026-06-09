package domain

import "time"

type Location struct {
	City      string
	Country   string
	Latitude  float64
	Longitude float64
	Timezone  string
}

type WeatherReport struct {
	Location  Location
	UpdatedAt time.Time

	Current CurrentWeather
	Metrics WeatherMetrics

	Daily  []DailyForecast
	Hourly []HourlyForecast
}

type CurrentWeather struct {
	TemperatureC     float64
	FeelsLikeC       float64
	Condition        string
	WeatherCode      int
	WindSpeedKmh     float64
	WindDirectionDeg int
}

type WeatherMetrics struct {
	HumidityPercent   int
	PressureHPa       float64
	PrecipitationMM   float64
	CloudCoverPercent int
	WindSpeedKmh      float64
	WindDirectionDeg  int
}

type DailyForecast struct {
	Date            time.Time
	MinTemperatureC float64
	MaxTemperatureC float64
	Condition       string
	WeatherCode     int
	PrecipitationMM float64
	MaxWindKmh      float64
}

type HourlyForecast struct {
	Time              time.Time
	TemperatureC      float64
	FeelsLikeC        float64
	Condition         string
	WeatherCode       int
	HumidityPercent   int
	PrecipitationMM   float64
	CloudCoverPercent int
	WindSpeedKmh      float64
}
