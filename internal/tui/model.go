package tui

import (
	"charm.land/bubbles/v2/help"

	"github.com/nfmdev/meteo/internal/domain"
)

type Model struct {
	report domain.WeatherReport

	selectedDay int
	showHelp    bool

	keys KeyMap
	help help.Model

	width  int
	height int
}

func NewModel(report domain.WeatherReport) Model {
	return Model{
		report:      report,
		selectedDay: 0,
		showHelp:    false,
		keys:        DefaultKeyMap(),
		help:        help.New(),
	}
}
