package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func DefaultForecastCacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("resolve user cache directory: %w", err)
	}

	return filepath.Join(cacheDir, AppName, ForecastCacheDir), nil
}

func ResolveForecastCacheDir(customDir string) (string, error) {
	customDir = strings.TrimSpace(customDir)
	if customDir == "" {
		return DefaultForecastCacheDir()
	}

	expandedDir, err := expandHomePath(customDir)
	if err != nil {
		return "", err
	}

	return filepath.Clean(expandedDir), nil
}

func ForecastCacheFilePath(cacheDir string, city string, country string) (string, error) {
	key, err := ForecastCacheKey(city, country)
	if err != nil {
		return "", err
	}

	return CacheFilePath(cacheDir, key)
}

func CacheFilePath(cacheDir string, key string) (string, error) {
	cacheDir = strings.TrimSpace(cacheDir)
	if cacheDir == "" {
		return "", ErrCacheDirRequired
	}
	if key == "" {
		return "", ErrCacheKeyRequired
	}

	return filepath.Join(cacheDir, key+CacheFileExtension), nil
}

func expandHomePath(path string) (string, error) {
	if path == "~" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve user home dir: %w", err)
		}

		return homeDir, nil
	}

	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("resolve user home dir: %w", err)
		}

		return filepath.Join(homeDir, strings.Trim(path, "~/")), nil
	}

	return path, nil
}
