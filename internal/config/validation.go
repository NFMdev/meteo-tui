package config

import (
	"fmt"
	"strings"

	"github.com/nfmdev/meteo/internal/domain"
)

func NormalizeConfig(config AppConfig) AppConfig {
	normalizedFavorites := make([]domain.SavedLocation, 0, len(config.Favorites))

	for _, favorite := range config.Favorites {
		normalizedFavorites = append(normalizedFavorites, domain.NormalizeSavedLocation(favorite))
	}

	return AppConfig{
		DefaultCity:    strings.TrimSpace(config.DefaultCity),
		DefaultCountry: strings.ToUpper(strings.TrimSpace(config.DefaultCountry)),
		Favorites:      normalizedFavorites,
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

	if err := validateFavorites(config.Favorites); err != nil {
		return err
	}

	return nil
}

func validateFavorites(favorites []domain.SavedLocation) error {
	seen := make(map[string]struct{}, len(favorites))

	for index, favorite := range favorites {
		if strings.TrimSpace(favorite.Name) == "" {
			return fmt.Errorf("%w at index %d", ErrFavoriteNameRequired, index)
		}
		if strings.TrimSpace(favorite.CountryCode) == "" {
			return fmt.Errorf("%w at index %d", ErrFavoriteCountryCodeRequired, index)
		}
		if !isValidCountryCode(favorite.CountryCode) {
			return fmt.Errorf("%w at favorite index %d", ErrInvalidCountryCode, index)
		}

		key := domain.SavedLocationKey(favorite)
		if key == "" {
			return fmt.Errorf("%w at index %d", ErrFavoriteNameRequired, index)
		}

		if _, exists := seen[key]; exists {
			return fmt.Errorf("%w: %s", ErrDuplicateFavorite, key)
		}

		seen[key] = struct{}{}
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
