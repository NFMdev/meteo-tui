package domain

import (
	"strings"
	"time"
)

const (
	WeatherProviderUnknown   = "Unknown"
	WeatherProviderOpenMeteo = "Open-Meteo"
)

type WeatherSource struct {
	Provider string    `json:"provider"`
	Cached   bool      `json:"cached"`
	CachedAt time.Time `json:"cached_at,omitempty"`
}

func NewFreshWeatherSource(provider string) WeatherSource {
	provider = strings.TrimSpace(provider)
	if provider == "" {
		provider = WeatherProviderUnknown
	}

	return WeatherSource{
		Provider: provider,
		Cached:   false,
	}
}

func NewCachedWeatherSource(provider string, cachedAt time.Time) WeatherSource {
	provider = strings.TrimSpace(provider)
	if provider == "" {
		provider = WeatherProviderUnknown
	}

	return WeatherSource{
		Provider: provider,
		Cached:   true,
		CachedAt: cachedAt,
	}
}
