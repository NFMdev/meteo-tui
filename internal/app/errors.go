package app

import "errors"

var (
	ErrCacheStoreRequired          = errors.New("cache store is required")
	ErrLocationSearcherRequired    = errors.New("location searcher is required")
	ErrLocationSearchQueryRequired = errors.New("location search query is required")
)
