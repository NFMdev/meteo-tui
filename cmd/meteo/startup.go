package main

import (
	"errors"
	"fmt"
	"strings"

	forecastCache "github.com/nfmdev/meteo/internal/cache"
	meteoConfig "github.com/nfmdev/meteo/internal/config"
)

const (
	fallbackCity    = "Copenhagen"
	fallbackCountry = "DK"
)

type startupOptions struct {
	city       string
	country    string
	configPath string
	cacheDir   string
	initConfig bool
	offline    bool
	fail       bool
}

type resolvedStartupOptions struct {
	city       string
	country    string
	configPath string
	cacheDir   string
	initConfig bool
	offline    bool
	fail       bool
}

func resolveStartupOptions(options startupOptions) (resolvedStartupOptions, error) {
	configPath, err := meteoConfig.ResolveConfigPath(options.configPath)
	if err != nil {
		return resolvedStartupOptions{}, fmt.Errorf("resolve config path: %w", err)
	}

	cacheDir, err := forecastCache.ResolveForecastCacheDir(options.cacheDir)
	if err != nil {
		return resolvedStartupOptions{}, fmt.Errorf("resolve cache directory: %w", err)
	}

	loadedConfig, err := loadOptionalConfig(configPath)
	if err != nil {
		return resolvedStartupOptions{}, err
	}

	city := strings.TrimSpace(options.city)
	if city == "" {
		city = loadedConfig.DefaultCity
	}
	if city == "" {
		city = fallbackCity
	}

	country := strings.TrimSpace(options.country)
	if country == "" {
		country = loadedConfig.DefaultCountry
	}
	if country == "" {
		country = fallbackCountry
	}

	finalConfig := meteoConfig.NormalizeConfig(meteoConfig.AppConfig{
		DefaultCity:    city,
		DefaultCountry: country,
	})

	if err := meteoConfig.ValidateConfig(finalConfig); err != nil {
		return resolvedStartupOptions{}, fmt.Errorf("validate startup location: %w", err)
	}

	return resolvedStartupOptions{
		city:       finalConfig.DefaultCity,
		country:    finalConfig.DefaultCountry,
		configPath: configPath,
		cacheDir:   cacheDir,
		initConfig: options.initConfig,
		offline:    options.offline,
		fail:       options.fail,
	}, nil
}

func loadOptionalConfig(path string) (meteoConfig.AppConfig, error) {
	config, err := meteoConfig.LoadConfig(path)
	if err == nil {
		return config, nil
	}

	if errors.Is(err, meteoConfig.ErrConfigNotFound) {
		return meteoConfig.AppConfig{}, nil
	}

	return meteoConfig.AppConfig{}, fmt.Errorf("load config: %w", err)
}
