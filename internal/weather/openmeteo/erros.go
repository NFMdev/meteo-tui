package openmeteo

import "errors"

var (
	ErrInvalidForecastResponse = errors.New("invalid forecast response")
	ErrForecastUnavailable     = errors.New("forecast data is unavailable")
)
