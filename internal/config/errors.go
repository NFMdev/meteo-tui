package config

import "errors"

var (
	ErrConfigPathRequired     = errors.New("config path required")
	ErrConfigNotFound         = errors.New("config file not found")
	ErrDefaultCityRequired    = errors.New("default city required")
	ErrDefaultCountryRequired = errors.New("default country required")
	ErrInvalidCountryCode     = errors.New("country code must be exactly two letters")
)
