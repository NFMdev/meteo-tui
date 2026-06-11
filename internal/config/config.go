package config

const (
	AppName        = "meteo-tui"
	ConfigFileName = "config.json"
)

type AppConfig struct {
	DefaultCity    string `json:"default_city"`
	DefaultCountry string `json:"default_country"`
}
