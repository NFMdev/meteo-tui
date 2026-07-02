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

type LocationPreferencesLoader interface {
	ListFavorites(ctx context.Context) ([]domain.SavedLocation, error)
	AddFavorite(ctx context.Context, location domain.SavedLocation) error
	RemoveFavorite(ctx context.Context, location domain.SavedLocation) error
	SetDefaultLocation(ctx context.Context, location domain.SavedLocation) error
}

type Model struct {
	city    string
	country string

	weatherLoader             WeatherLoader
	locationSearchLoader      LocationSearchLoader
	locationPreferencesLoader LocationPreferencesLoader

	report domain.WeatherReport

	loading bool
	err     error

	selectedDay int
	showHelp    bool
	keys        KeyMap
	help        help.Model
	spinner     spinner.Model
	viewport    viewport.Model

	mode screenMode

	searchInput          textinput.Model
	searching            bool
	searchResults        []domain.LocationSearchResult
	selectedSearchResult int
	searchErr            error

	favorites        []domain.SavedLocation
	favoritesLoading bool
	selectedFavorite int
	favoritesErr     error

	statusMessage string

	width  int
	height int
}

func NewModel(
	city string,
	country string,
	weatherLoader WeatherLoader,
) Model {
	return NewModelWithSearchAndPreferences(city, country, weatherLoader, nil, nil)
}

func NewModelWithSearch(
	city string,
	country string,
	weatherLoader WeatherLoader,
	locationSearchLoader LocationSearchLoader,
) Model {
	return NewModelWithSearchAndPreferences(
		city,
		country,
		weatherLoader,
		locationSearchLoader,
		nil,
	)
}

func NewModelWithSearchAndPreferences(
	city string,
	country string,
	weatherLoader WeatherLoader,
	locationSearchLoader LocationSearchLoader,
	locationPreferencesLoader LocationPreferencesLoader,
) Model {
	return Model{
		city:                      city,
		country:                   country,
		weatherLoader:             weatherLoader,
		locationSearchLoader:      locationSearchLoader,
		locationPreferencesLoader: locationPreferencesLoader,
		loading:                   true,
		selectedDay:               0,
		showHelp:                  false,
		keys:                      DefaultKeyMap(),
		help:                      help.New(),
		spinner:                   spinner.New(spinner.WithSpinner(spinner.Dot)),
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
