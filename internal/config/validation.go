package config

import "strings"

func NormalizeConfig(config AppConfig) AppConfig {
	return AppConfig{
		DefaultCity:    strings.TrimSpace(config.DefaultCity),
		DefaultCountry: strings.ToUpper(strings.TrimSpace(config.DefaultCountry)),
	}
}

func ValidateConfig(config AppConfig) error {
	config = NormalizeConfig(config)

	if config.DefaultCity == "" {
		return ErrDefaultCityRequired
	}

	if config.DefaultCountry == "" {
		return ErrDefaultCountryRequired
	}

	if !isValidCountryCode(config.DefaultCountry) {
		return ErrInvalidCountryCode
	}

	return nil
}

func isValidCountryCode(country string) bool {
	if len(country) != 2 {
		return false
	}
	for _, char := range country {
		if char < 'A' || char > 'Z' {
			return false
		}
	}
	return true
}
