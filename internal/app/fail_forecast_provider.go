package app

import (
	"context"
	"errors"

	"github.com/nfmdev/meteo/internal/domain"
)

var errSimulatedForecastFailure = errors.New("simulated forecast provider failure")

type FailForecastProvider struct{}

func (p FailForecastProvider) GetForecast(
	ctx context.Context,
	location domain.Location,
) (domain.WeatherReport, error) {
	return domain.WeatherReport{}, errSimulatedForecastFailure
}
