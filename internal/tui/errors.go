package tui

import "errors"

var (
	errSearchQueryRequired               = errors.New("search query is required")
	errLocationSearchLoaderRequired      = errors.New("location search loader is required")
	errSearchResultRequired              = errors.New("search result is required")
	errLocationPreferencesLoaderRequired = errors.New("locatoin preferences loader is required")
)
