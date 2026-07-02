package config

import "github.com/nfmdev/meteo/internal/domain"

const (
	AppName        = "meteo-tui"
	ConfigFileName = "config.json"
)

type AppConfig struct {
	DefaultCity    string                 `json:"default_city"`
	DefaultCountry string                 `json:"default_country"`
	Favorites      []domain.SavedLocation `json:"favorites,omitempty"`
}
