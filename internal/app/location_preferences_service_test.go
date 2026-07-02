package app

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	meteoConfig "github.com/nfmdev/meteo/internal/config"
	"github.com/nfmdev/meteo/internal/domain"
)

func TestNewConfigLocationPreferencesServiceRejectsEmptyPath(t *testing.T) {
	t.Parallel()

	_, err := NewConfigLocationPreferencesService("   ")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrLocationPreferencesConfigPathRequired) {
		t.Fatalf("expected ErrLocationPreferencesConfigPathRequired, got %v", err)
	}
}

func TestLocationPreferencesServiceListFavoritesReturnsEmptyWhenConfigMissing(t *testing.T) {
	t.Parallel()

	service := newTestLocationPreferencesService(t, missingTestConfigPath(t))

	favorites, err := service.ListFavorites(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if favorites == nil {
		t.Fatal("expected empty slice, got nil")
	}

	if len(favorites) != 0 {
		t.Fatalf("expected 0 favorites, got %d", len(favorites))
	}
}

func TestLocationPreferencesServiceListFavoritesReturnsCopy(t *testing.T) {
	t.Parallel()

	path := writePreferencesTestConfig(t, meteoConfig.AppConfig{
		DefaultCity:    "Aalborg",
		DefaultCountry: "DK",
		Favorites: []domain.SavedLocation{
			madridSavedLocation(),
		},
	})

	service := newTestLocationPreferencesService(t, path)

	favorites, err := service.ListFavorites(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	favorites[0].Name = "Mutated"

	reloaded, err := service.ListFavorites(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if reloaded[0].Name != "Madrid" {
		t.Fatalf("expected stored favorite to remain Madrid, got %q", reloaded[0].Name)
	}
}

func TestLocationPreferencesServiceAddFavoriteCreatesConfigWhenMissing(t *testing.T) {
	t.Parallel()

	path := missingTestConfigPath(t)

	service := newTestLocationPreferencesService(t, path)

	err := service.AddFavorite(context.Background(), madridSavedLocation())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	config, err := meteoConfig.LoadConfig(path)
	if err != nil {
		t.Fatalf("expected config to be created, got %v", err)
	}

	if config.DefaultCity != "Madrid" {
		t.Fatalf("expected default city Madrid, got %q", config.DefaultCity)
	}

	if config.DefaultCountry != "ES" {
		t.Fatalf("expected default country ES, got %q", config.DefaultCountry)
	}

	if len(config.Favorites) != 1 {
		t.Fatalf("expected 1 favorite, got %d", len(config.Favorites))
	}

	if config.Favorites[0].Name != "Madrid" {
		t.Fatalf("expected favorite Madrid, got %q", config.Favorites[0].Name)
	}
}

func TestLocationPreferencesServiceAddFavoritePersistsFavorite(t *testing.T) {
	t.Parallel()

	path := writePreferencesTestConfig(t, meteoConfig.AppConfig{
		DefaultCity:    "Aalborg",
		DefaultCountry: "DK",
	})

	service := newTestLocationPreferencesService(t, path)

	err := service.AddFavorite(context.Background(), madridSavedLocation())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	config, err := meteoConfig.LoadConfig(path)
	if err != nil {
		t.Fatalf("expected no error loading config, got %v", err)
	}

	if len(config.Favorites) != 1 {
		t.Fatalf("expected 1 favorite, got %d", len(config.Favorites))
	}

	favorite := config.Favorites[0]

	if favorite.Name != "Madrid" {
		t.Fatalf("expected Madrid, got %q", favorite.Name)
	}

	if favorite.CountryCode != "ES" {
		t.Fatalf("expected ES, got %q", favorite.CountryCode)
	}
}

func TestLocationPreferencesServiceAddFavoriteNormalizesFavorite(t *testing.T) {
	t.Parallel()

	path := writePreferencesTestConfig(t, meteoConfig.AppConfig{
		DefaultCity:    "Aalborg",
		DefaultCountry: "DK",
	})

	service := newTestLocationPreferencesService(t, path)

	err := service.AddFavorite(context.Background(), domain.SavedLocation{
		Name:        " Madrid ",
		Country:     " Spain ",
		CountryCode: " es ",
		Admin1:      " Community of Madrid ",
		Timezone:    " Europe/Madrid ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	config, err := meteoConfig.LoadConfig(path)
	if err != nil {
		t.Fatalf("expected no error loading config, got %v", err)
	}

	favorite := config.Favorites[0]

	if favorite.Name != "Madrid" {
		t.Fatalf("expected Madrid, got %q", favorite.Name)
	}

	if favorite.CountryCode != "ES" {
		t.Fatalf("expected ES, got %q", favorite.CountryCode)
	}

	if favorite.Admin1 != "Community of Madrid" {
		t.Fatalf("expected Community of Madrid, got %q", favorite.Admin1)
	}
}

func TestLocationPreferencesServiceAddFavoriteIsIdempotent(t *testing.T) {
	t.Parallel()

	path := writePreferencesTestConfig(t, meteoConfig.AppConfig{
		DefaultCity:    "Aalborg",
		DefaultCountry: "DK",
		Favorites: []domain.SavedLocation{
			madridSavedLocation(),
		},
	})

	service := newTestLocationPreferencesService(t, path)

	err := service.AddFavorite(context.Background(), domain.SavedLocation{
		Name:        " Madrid ",
		CountryCode: " es ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	config, err := meteoConfig.LoadConfig(path)
	if err != nil {
		t.Fatalf("expected no error loading config, got %v", err)
	}

	if len(config.Favorites) != 1 {
		t.Fatalf("expected 1 favorite, got %d", len(config.Favorites))
	}
}

func TestLocationPreferencesServiceAddFavoriteRejectsInvalidLocation(t *testing.T) {
	t.Parallel()

	service := newTestLocationPreferencesService(t, missingTestConfigPath(t))

	err := service.AddFavorite(context.Background(), domain.SavedLocation{
		Name:        "Madrid",
		CountryCode: "ESP",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, meteoConfig.ErrInvalidCountryCode) {
		t.Fatalf("expected ErrInvalidCountryCode, got %v", err)
	}
}

func TestLocationPreferencesServiceRemoveFavorite(t *testing.T) {
	t.Parallel()

	path := writePreferencesTestConfig(t, meteoConfig.AppConfig{
		DefaultCity:    "Aalborg",
		DefaultCountry: "DK",
		Favorites: []domain.SavedLocation{
			madridSavedLocation(),
			copenhagenSavedLocation(),
		},
	})

	service := newTestLocationPreferencesService(t, path)

	err := service.RemoveFavorite(context.Background(), domain.SavedLocation{
		Name:        " Madrid ",
		CountryCode: " es ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	config, err := meteoConfig.LoadConfig(path)
	if err != nil {
		t.Fatalf("expected no error loading config, got %v", err)
	}

	if len(config.Favorites) != 1 {
		t.Fatalf("expected 1 favorite, got %d", len(config.Favorites))
	}

	if config.Favorites[0].Name != "Copenhagen" {
		t.Fatalf("expected Copenhagen to remain, got %q", config.Favorites[0].Name)
	}
}

func TestLocationPreferencesServiceRemoveFavoriteReturnsNotFound(t *testing.T) {
	t.Parallel()

	path := writePreferencesTestConfig(t, meteoConfig.AppConfig{
		DefaultCity:    "Aalborg",
		DefaultCountry: "DK",
		Favorites: []domain.SavedLocation{
			copenhagenSavedLocation(),
		},
	})

	service := newTestLocationPreferencesService(t, path)

	err := service.RemoveFavorite(context.Background(), madridSavedLocation())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrFavoriteNotFound) {
		t.Fatalf("expected ErrFavoriteNotFound, got %v", err)
	}
}

func TestLocationPreferencesServiceRemoveFavoriteReturnsNotFoundWhenConfigMissing(t *testing.T) {
	t.Parallel()

	service := newTestLocationPreferencesService(t, missingTestConfigPath(t))

	err := service.RemoveFavorite(context.Background(), madridSavedLocation())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, ErrFavoriteNotFound) {
		t.Fatalf("expected ErrFavoriteNotFound, got %v", err)
	}
}

func TestLocationPreferencesServiceSetDefaultLocation(t *testing.T) {
	t.Parallel()

	path := writePreferencesTestConfig(t, meteoConfig.AppConfig{
		DefaultCity:    "Aalborg",
		DefaultCountry: "DK",
		Favorites: []domain.SavedLocation{
			madridSavedLocation(),
		},
	})

	service := newTestLocationPreferencesService(t, path)

	err := service.SetDefaultLocation(context.Background(), madridSavedLocation())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	config, err := meteoConfig.LoadConfig(path)
	if err != nil {
		t.Fatalf("expected no error loading config, got %v", err)
	}

	if config.DefaultCity != "Madrid" {
		t.Fatalf("expected default city Madrid, got %q", config.DefaultCity)
	}

	if config.DefaultCountry != "ES" {
		t.Fatalf("expected default country ES, got %q", config.DefaultCountry)
	}

	if len(config.Favorites) != 1 {
		t.Fatalf("expected favorite to be preserved, got %d", len(config.Favorites))
	}
}

func TestLocationPreferencesServiceSetDefaultLocationCreatesConfigWhenMissing(t *testing.T) {
	t.Parallel()

	path := missingTestConfigPath(t)

	service := newTestLocationPreferencesService(t, path)

	err := service.SetDefaultLocation(context.Background(), madridSavedLocation())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	config, err := meteoConfig.LoadConfig(path)
	if err != nil {
		t.Fatalf("expected config to be created, got %v", err)
	}

	if config.DefaultCity != "Madrid" {
		t.Fatalf("expected default city Madrid, got %q", config.DefaultCity)
	}

	if config.DefaultCountry != "ES" {
		t.Fatalf("expected default country ES, got %q", config.DefaultCountry)
	}

	if len(config.Favorites) != 0 {
		t.Fatalf("expected no favorites, got %d", len(config.Favorites))
	}
}

func newTestLocationPreferencesService(
	t *testing.T,
	configPath string,
) ConfigLocationPreferencesService {
	t.Helper()

	service, err := NewConfigLocationPreferencesService(configPath)
	if err != nil {
		t.Fatalf("expected no error creating service, got %v", err)
	}

	return service
}

func missingTestConfigPath(t *testing.T) string {
	t.Helper()

	return filepath.Join(t.TempDir(), "config.json")
}

func writePreferencesTestConfig(
	t *testing.T,
	config meteoConfig.AppConfig,
) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.json")

	if err := meteoConfig.WriteConfig(path, config); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	return path
}

func madridSavedLocation() domain.SavedLocation {
	return domain.SavedLocation{
		Name:        "Madrid",
		Country:     "Spain",
		CountryCode: "ES",
		Admin1:      "Community of Madrid",
		Latitude:    40.4168,
		Longitude:   -3.7038,
		Timezone:    "Europe/Madrid",
	}
}

func copenhagenSavedLocation() domain.SavedLocation {
	return domain.SavedLocation{
		Name:        "Copenhagen",
		Country:     "Denmark",
		CountryCode: "DK",
		Admin1:      "Capital Region",
		Latitude:    55.6761,
		Longitude:   12.5683,
		Timezone:    "Europe/Copenhagen",
	}
}
