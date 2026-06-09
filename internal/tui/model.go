package tui

import (
	"context"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/spinner"

	"github.com/nfmdev/meteo/internal/domain"
)

type WeatherLoader interface {
	GetWeather(ctx context.Context, city string, country string) (domain.WeatherReport, error)
}

type Model struct {
	city    string
	country string

	loader WeatherLoader

	report domain.WeatherReport

	loading bool
	err     error

	selectedDay int
	showHelp    bool

	keys    KeyMap
	help    help.Model
	spinner spinner.Model

	width  int
	height int
}

func NewModel(city string, country string, loader WeatherLoader) Model {
	return Model{
		city:        city,
		country:     country,
		loader:      loader,
		loading:     true,
		selectedDay: 0,
		showHelp:    false,
		keys:        DefaultKeyMap(),
		help:        help.New(),
		spinner:     spinner.New(spinner.WithSpinner(spinner.Dot)),
	}
}
