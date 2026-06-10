package location

import "errors"

var (
	ErrCityRequired         = errors.New("city is required")
	ErrCountryRequired      = errors.New("country is required")
	ErrLocationNotFound     = errors.New("location not found")
	ErrInvalidGeocoodeReply = errors.New("invalid geocoding response")
)
