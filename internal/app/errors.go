package app

import "errors"

var (
	ErrCacheStoreRequired                    = errors.New("cache store is required")
	ErrLocationSearcherRequired              = errors.New("location searcher is required")
	ErrLocationSearchQueryRequired           = errors.New("location search query is required")
	ErrLocationPreferencesConfigPathRequired = errors.New("locaiton preferences config path is required")
	ErrSavedLocationNameRequired             = errors.New("saved location name is required")
	ErrSavedLocationCountryCodeRequired      = errors.New("saved location country code is required")
	ErrFavoriteNotFound                      = errors.New("favorite location not found")
)
