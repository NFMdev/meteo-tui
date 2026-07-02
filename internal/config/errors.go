package config

import "errors"

var (
	ErrConfigPathRequired          = errors.New("config path required")
	ErrConfigNotFound              = errors.New("config file not found")
	ErrDefaultCityRequired         = errors.New("default city required")
	ErrDefaultCountryRequired      = errors.New("default country required")
	ErrInvalidCountryCode          = errors.New("country code must be exactly two letters")
	ErrFavoriteNameRequired        = errors.New("favorite name is required")
	ErrFavoriteCountryCodeRequired = errors.New("favorite country code is required")
	ErrDuplicateFavorite           = errors.New("duplicate favorite location")
)
