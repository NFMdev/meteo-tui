package tui

import "github.com/nfmdev/meteo/internal/domain"

type weatherLoadedMsg struct {
	report domain.WeatherReport
}

type weatherFailedMsg struct {
	err error
}

type locationSearchLoadedMsg struct {
	results []domain.LocationSearchResult
}

type locationSearchFailedMsg struct {
	err error
}

type favoritesLoadedMsg struct {
	favorites []domain.SavedLocation
}

type favoritesFailedMsg struct {
	err error
}

type locationPreferenceUpdatedMsg struct {
	message string
}

type locationPreferenceFailedMsg struct {
	err error
}
