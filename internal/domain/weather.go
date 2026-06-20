package domain

import "time"

type Location struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"logitude"`
	Timezone  string  `json:"timezone"`
}

type WeatherReport struct {
	Location  Location         `json:"location"`
	UpdatedAt time.Time        `json:"updated_at"`
	Source    WeatherSource    `json:"weather_source"`
	Current   CurrentWeather   `json:"current"`
	Metrics   WeatherMetrics   `json:"metrics"`
	Daily     []DailyForecast  `json:"daily"`
	Hourly    []HourlyForecast `json:"hourly"`
}

type CurrentWeather struct {
	TemperatureC     float64 `json:"temperature_c"`
	FeelsLikeC       float64 `json:"feels_like_c"`
	Condition        string  `json:"condition"`
	WeatherCode      int     `json:"weather_code"`
	WindSpeedKmh     float64 `json:"wind_speed_kmh"`
	WindDirectionDeg int     `json:"wind_direction_deg"`
}

type WeatherMetrics struct {
	HumidityPercent   int     `json:"humidity_percent"`
	PressureHPa       float64 `json:"pressure_hpa"`
	PrecipitationMM   float64 `json:"precipitation_mm"`
	CloudCoverPercent int     `json:"cloud_cover_percent"`
	WindSpeedKmh      float64 `json:"wind_speed_kmh"`
	WindDirectionDeg  int     `json:"wind_direction_deg"`
}

type DailyForecast struct {
	Date            time.Time `json:"date"`
	MinTemperatureC float64   `json:"min_temperature_c"`
	MaxTemperatureC float64   `json:"max_temperature_c"`
	Condition       string    `json:"condition"`
	WeatherCode     int       `json:"weather_code"`
	PrecipitationMM float64   `json:"precipitation_mm"`
	MaxWindKmh      float64   `json:"max_wind_kmh"`
}

type HourlyForecast struct {
	Time              time.Time `json:"time"`
	TemperatureC      float64   `json:"temperature_c"`
	FeelsLikeC        float64   `json:"feels_like_c"`
	Condition         string    `json:"condition"`
	WeatherCode       int       `json:"weather_code"`
	HumidityPercent   int       `json:"humidity_percent"`
	PrecipitationMM   float64   `json:"precipitation_mm"`
	CloudCoverPercent int       `json:"cloud_cover_percent"`
	WindSpeedKmh      float64   `json:"wind_speed_kmh"`
}
