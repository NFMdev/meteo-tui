package tui

import "github.com/nfmdev/meteo/internal/domain"

type weatherLoadedMsg struct {
	report domain.WeatherReport
}

type weatherFailedMsg struct {
	err error
}
