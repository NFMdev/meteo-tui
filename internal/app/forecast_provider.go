package app

import (
	"context"

	"github.com/nfmdev/meteo/internal/domain"
)

type ForecastProvider interface {
	GetForecast(ctx context.Context, location domain.Location) (domain.WeatherReport, error)
}
