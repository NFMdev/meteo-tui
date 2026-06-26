package tui

import (
	"context"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"

	"github.com/nfmdev/meteo/internal/domain"
)

type WeatherLoader interface {
	GetWeather(ctx context.Context, city string, country string) (domain.WeatherReport, error)
}

type LocationSearchLoader interface {
	SearchLocations(ctx context.Context, query string) ([]domain.LocationSearchResult, error)
}

type Model struct {
	city    string
	country string

	loader               WeatherLoader
	locationSearchLoader LocationSearchLoader

	report domain.WeatherReport

	loading bool
	err     error

	selectedDay int
	showHelp    bool
	keys        KeyMap
	help        help.Model
	spinner     spinner.Model
	viewport    viewport.Model

	mode                 screenMode
	searchInput          textinput.Model
	searching            bool
	searchResults        []domain.LocationSearchResult
	selectedSearchResult int
	searchErr            error

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
		viewport: viewport.New(
			viewport.WithWidth(defaultTerminalWidth),
			viewport.WithHeight(20),
		),
		mode:        screenModeDashboard,
		searchInput: newSearchInput(),
	}
}

func newSearchInput() textinput.Model {
	input := textinput.New()
	input.Placeholder = "Search city..."
	input.Prompt = "> "
	input.CharLimit = 80
	input.SetWidth(40)

	return input
}

func NewModelWithSearch(
	city string,
	country string,
	loader WeatherLoader,
	locationSearchLoader LocationSearchLoader,
) Model {
	return Model{
		city:                 city,
		country:              country,
		loader:               loader,
		locationSearchLoader: locationSearchLoader,
		loading:              true,
		selectedDay:          0,
		showHelp:             false,
		keys:                 DefaultKeyMap(),
		help:                 help.New(),
		spinner:              spinner.New(spinner.WithSpinner(spinner.Dot)),
		viewport: viewport.New(
			viewport.WithWidth(defaultTerminalWidth),
			viewport.WithHeight(20),
		),
		mode:        screenModeDashboard,
		searchInput: newSearchInput(),
	}
}
