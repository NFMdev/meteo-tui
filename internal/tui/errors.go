package tui

import "errors"

var (
	errSearchQueryRequired               = errors.New("search query is required")
	errLocationSearchLoaderRequired      = errors.New("location search loader is required")
	errSearchResultRequired              = errors.New("search result is required")
	errLocationPreferencesLoaderRequired = errors.New("location preferences loader is required")
	errCurrentLocationRequired           = errors.New("current location is required")
)
