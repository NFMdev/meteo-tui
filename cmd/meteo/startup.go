package main

import (
	"errors"
	"fmt"
	"strings"

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
	initConfig bool
	fail       bool
	fake       bool
}

type resolvedStartupOptions struct {
	city       string
	country    string
	configPath string
	initConfig bool
	fail       bool
	fake       bool
}

func resolveStartupOptions(options startupOptions) (resolvedStartupOptions, error) {
	configPath, err := meteoConfig.ResolveConfigPath(options.configPath)
	if err != nil {
		return resolvedStartupOptions{}, fmt.Errorf("resolve config path: %w", err)
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
		initConfig: options.initConfig,
		fail:       options.fail,
		fake:       options.fake,
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
