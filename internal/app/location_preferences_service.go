package app

import (
	"context"
	"errors"
	"fmt"
	"strings"

	meteoConfig "github.com/nfmdev/meteo/internal/config"
	"github.com/nfmdev/meteo/internal/domain"
)

type LocationPreferencesService interface {
	ListFavorites(ctx context.Context) ([]domain.SavedLocation, error)
	AddFavorite(ctx context.Context, location domain.SavedLocation) error
	RemoveFavorite(ctx context.Context, location domain.SavedLocation) error
	SetDefaultLocation(ctx context.Context, location domain.SavedLocation) error
}

type ConfigLocationPreferencesService struct {
	configPath string
}

func NewConfigLocationPreferencesService(
	configPath string,
) (ConfigLocationPreferencesService, error) {
	configPath = strings.TrimSpace(configPath)
	if configPath == "" {
		return ConfigLocationPreferencesService{},
			ErrLocationPreferencesConfigPathRequired
	}

	return ConfigLocationPreferencesService{
		configPath: configPath,
	}, nil
}

func (s ConfigLocationPreferencesService) ListFavorites(
	ctx context.Context,
) ([]domain.SavedLocation, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	config, err := meteoConfig.LoadConfig(s.configPath)
	if err != nil {
		if errors.Is(err, meteoConfig.ErrConfigNotFound) {
			return []domain.SavedLocation{}, nil
		}

		return nil, err
	}

	favorites := make([]domain.SavedLocation, len(config.Favorites))
	copy(favorites, config.Favorites)

	return favorites, nil
}

func (s ConfigLocationPreferencesService) AddFavorite(
	ctx context.Context,
	location domain.SavedLocation,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	location, err := validateSavedLocation(location)
	if err != nil {
		return err
	}

	config, err := s.loadConfigForMutation(location)
	if err != nil {
		return err
	}

	if containsSavedLocation(config.Favorites, location) {
		return nil
	}

	config.Favorites = append(config.Favorites, location)

	return meteoConfig.WriteConfig(s.configPath, config)
}

func (s ConfigLocationPreferencesService) RemoveFavorite(
	ctx context.Context,
	location domain.SavedLocation,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	location, err := validateSavedLocation(location)
	if err != nil {
		return err
	}

	config, err := meteoConfig.LoadConfig(s.configPath)
	if err != nil {
		if errors.Is(err, meteoConfig.ErrConfigNotFound) {
			return ErrFavoriteNotFound
		}

		return err
	}

	updatedFavorites := make([]domain.SavedLocation, 0, len(config.Favorites))
	removed := false

	for _, favorite := range config.Favorites {
		if domain.SameSavedLocation(favorite, location) {
			removed = true
			continue
		}

		updatedFavorites = append(updatedFavorites, favorite)
	}

	if !removed {
		return ErrFavoriteNotFound
	}

	config.Favorites = updatedFavorites

	return meteoConfig.WriteConfig(s.configPath, config)
}

func (s ConfigLocationPreferencesService) SetDefaultLocation(
	ctx context.Context,
	location domain.SavedLocation,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	location, err := validateSavedLocation(location)
	if err != nil {
		return err
	}

	config, err := s.loadConfigForMutation(location)
	if err != nil {
		return err
	}

	config.DefaultCity = location.Name
	config.DefaultCountry = location.CountryCode

	return meteoConfig.WriteConfig(s.configPath, config)
}

func (s ConfigLocationPreferencesService) loadConfigForMutation(
	fallbackLocation domain.SavedLocation,
) (meteoConfig.AppConfig, error) {
	config, err := meteoConfig.LoadConfig(s.configPath)
	if err == nil {
		return config, nil
	}

	if !errors.Is(err, meteoConfig.ErrConfigNotFound) {
		return meteoConfig.AppConfig{}, err
	}

	return meteoConfig.AppConfig{
		DefaultCity:    fallbackLocation.Name,
		DefaultCountry: fallbackLocation.CountryCode,
		Favorites:      []domain.SavedLocation{},
	}, nil
}

func validateSavedLocation(
	location domain.SavedLocation,
) (domain.SavedLocation, error) {
	location = domain.NormalizeSavedLocation(location)

	if location.Name == "" {
		return domain.SavedLocation{}, ErrSavedLocationNameRequired
	}

	if location.CountryCode == "" {
		return domain.SavedLocation{}, ErrSavedLocationCountryCodeRequired
	}

	if len(location.CountryCode) != 2 {
		return domain.SavedLocation{}, fmt.Errorf(
			"%w: %s",
			meteoConfig.ErrInvalidCountryCode,
			location.CountryCode,
		)
	}

	for _, char := range location.CountryCode {
		if char < 'A' || char > 'Z' {
			return domain.SavedLocation{}, fmt.Errorf(
				"%w: %s",
				meteoConfig.ErrInvalidCountryCode,
				location.CountryCode,
			)
		}
	}

	return location, nil
}

func containsSavedLocation(
	favorites []domain.SavedLocation,
	location domain.SavedLocation,
) bool {
	for _, favorite := range favorites {
		if domain.SameSavedLocation(favorite, location) {
			return true
		}
	}

	return false
}
