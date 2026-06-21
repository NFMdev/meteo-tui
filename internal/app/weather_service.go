package app

import (
	"context"
	"errors"

	"github.com/nfmdev/meteo/internal/domain"
)

var ErrWeatherUnavailable = errors.New("weather data is currently unavailable")

type WeatherService interface {
	GetWeather(ctx context.Context, city string, country string) (domain.WeatherReport, error)
}

type RealWeatherService struct {
	locationResolver LocationResolver
	forecastProvider ForecastProvider
	cacheStore       ForecastCacheStore
	options          WeatherServiceOptions
}

type WeatherServiceOptions struct {
	Offline bool
}

// TODO: remove fail weather service before releasing v1.0
// Keeping it for development testing
type FailWeatherService struct {
	locationResolver LocationResolver
	forecastProvider ForecastProvider
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

func NewWeatherServiceWithCache(
	locationResolver LocationResolver,
	forecastProvider ForecastProvider,
	cacheStore ForecastCacheStore,
	options WeatherServiceOptions,
) RealWeatherService {
	return RealWeatherService{
		locationResolver: locationResolver,
		forecastProvider: forecastProvider,
		cacheStore:       cacheStore,
		options:          options,
	}
}

func NewFailingWeatherService(locationResolver LocationResolver, forecastProvider ForecastProvider) FailWeatherService {
	return FailWeatherService{
		locationResolver: locationResolver,
		forecastProvider: forecastProvider,
		shouldFail:       true,
	}
}

func (s RealWeatherService) GetWeather(
	ctx context.Context,
	city string,
	country string,
) (domain.WeatherReport, error) {
	if s.options.Offline {
		return s.getCachedWeather(city, country)
	}

	location, err := s.locationResolver.Resolve(ctx, city, country)
	if err != nil {
		return s.getCachedWeatherOrOriginalError(city, country, err)
	}

	report, err := s.forecastProvider.GetForecast(ctx, location)
	if err != nil {
		return s.getCachedWeatherOrOriginalError(city, country, err)
	}
	s.writeCache(city, country, report)

	return report, nil
}

func (s RealWeatherService) getCachedWeather(
	city string,
	country string,
) (domain.WeatherReport, error) {
	if s.cacheStore == nil {
		return domain.WeatherReport{}, ErrCacheStoreRequired
	}

	return s.cacheStore.ReadReport(city, country)
}

func (s RealWeatherService) getCachedWeatherOrOriginalError(
	city string,
	country string,
	originalErr error,
) (domain.WeatherReport, error) {
	if s.cacheStore == nil {
		return domain.WeatherReport{}, originalErr
	}

	report, err := s.cacheStore.ReadReport(city, country)
	if err != nil {
		return domain.WeatherReport{}, originalErr
	}

	return report, nil
}

func (s RealWeatherService) writeCache(
	city string,
	country string,
	report domain.WeatherReport,
) {
	if s.cacheStore == nil {
		return
	}
	s.cacheStore.WriteForecast(city, country, report)
}
