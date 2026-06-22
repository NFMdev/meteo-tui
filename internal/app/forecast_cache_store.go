package app

import "github.com/nfmdev/meteo/internal/domain"

type ForecastCacheStore interface {
	WriteForecast(city string, country string, report domain.WeatherReport) error
	ReadReport(city string, country string) (domain.WeatherReport, error)
}
