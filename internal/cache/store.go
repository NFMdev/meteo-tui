package cache

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nfmdev/meteo/internal/domain"
)

type ForecastStore struct {
	cacheDir string
	now      func() time.Time
}

func NewForecastStore(cacheDir string) (ForecastStore, error) {
	cacheDir = strings.TrimSpace(cacheDir)
	if cacheDir == "" {
		return ForecastStore{}, ErrCacheDirRequired
	}

	return ForecastStore{
		cacheDir: filepath.Clean(cacheDir),
		now:      time.Now,
	}, nil
}

func (s ForecastStore) WriteForecast(
	city string,
	country string,
	report domain.WeatherReport,
) error {
	key, err := ForecastCacheKey(city, country)
	if err != nil {
		return err
	}

	normalizedCity := strings.TrimSpace(city)
	normalizedCountry := strings.ToUpper(strings.TrimSpace(country))

	entry := ForecastCacheEntry{
		Key:      key,
		City:     normalizedCity,
		Country:  normalizedCountry,
		CachedAt: s.now(),
		Report:   report,
	}

	if err := validateForecastCacheEntry(entry, key); err != nil {
		return err
	}

	path, err := CacheFilePath(s.cacheDir, key)
	if err != nil {
		return err
	}

	parentDir := filepath.Dir(path)
	if err := os.MkdirAll(parentDir, 0o755); err != nil {
		return fmt.Errorf("create forecast cache directory %q: %w", parentDir, err)
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("encode forecast cache entry: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write forecast cache file %q: %w", path, err)
	}

	return nil
}

func (s ForecastStore) ReadForecast(
	city string,
	country string,
) (ForecastCacheEntry, error) {
	key, err := ForecastCacheKey(city, country)
	if err != nil {
		return ForecastCacheEntry{}, err
	}

	path, err := CacheFilePath(s.cacheDir, key)
	if err != nil {
		return ForecastCacheEntry{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ForecastCacheEntry{}, fmt.Errorf("%w: %s", ErrCacheNotFound, path)
		}

		return ForecastCacheEntry{}, fmt.Errorf("read forecast cache file %q: %w", path, err)
	}

	var entry ForecastCacheEntry

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&entry); err != nil {
		return ForecastCacheEntry{}, fmt.Errorf("decode forecast cache file %q: %w", path, err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return ForecastCacheEntry{}, fmt.Errorf("decode forecast cache file %q: multiple JSON values are not allowed", path)
	}
	if err := validateForecastCacheEntry(entry, key); err != nil {
		return ForecastCacheEntry{}, fmt.Errorf("validate forecast cache file %q: %w", path, err)
	}

	entry.Report.Source = domain.NewCachedWeatherSource(
		entry.Report.Source.Provider,
		entry.CachedAt,
	)

	return entry, nil
}

func validateForecastCacheEntry(entry ForecastCacheEntry, expectedKey string) error {
	if strings.TrimSpace(entry.Key) == "" {
		return fmt.Errorf("%w: missing key", ErrInvalidCacheEntry)
	}

	if expectedKey != "" && entry.Key != expectedKey {
		return fmt.Errorf(
			"%w: expected key %q, got %q",
			ErrInvalidCacheKey,
			expectedKey,
			entry.Key,
		)
	}

	if strings.TrimSpace(entry.City) == "" {
		return fmt.Errorf("%w: missing city", ErrInvalidCacheEntry)
	}

	if strings.TrimSpace(entry.Country) == "" {
		return fmt.Errorf("%w: missing country", ErrInvalidCacheEntry)
	}

	if entry.CachedAt.IsZero() {
		return fmt.Errorf("%w: missing cached_at", ErrInvalidCacheEntry)
	}

	return nil
}
