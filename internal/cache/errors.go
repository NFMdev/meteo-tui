package cache

import "errors"

var (
	ErrCacheDirRequired     = errors.New("cahceh directory is required")
	ErrCacheKeyRequired     = errors.New("cache key is required")
	ErrCacheCityRequired    = errors.New("cache city is required")
	ErrCacheCountryRequired = errors.New("cache country is required")
	ErrInvalidCacheCountry  = errors.New("cache contry must be exactly two letters")
	ErrInvalidCacheKey      = errors.New("cache key is invalid")

	ErrCacheNotFound     = errors.New("forecast cache not found")
	ErrInvalidCacheEntry = errors.New("invalid forecast cache entry")
)
