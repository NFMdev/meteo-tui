package cache

import (
	"time"

	"github.com/nfmdev/meteo/internal/domain"
)

const (
	AppName            = "Meteo"
	ForecastCacheDir   = "forecasts"
	CacheFileExtension = ".json"
)

type ForecastCacheEntry struct {
	Key      string               `json:"key"`
	City     string               `json:"city"`
	Country  string               `json:"country"`
	CachedAt time.Time            `json:"cached_at"`
	Report   domain.WeatherReport `json:"report"`
}
