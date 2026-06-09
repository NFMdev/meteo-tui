package tui

import "github.com/nfmdev/meteo/internal/domain"

type Model struct {
	report domain.WeatherReport

	width  int
	height int
}

func NewModel(report domain.WeatherReport) Model {
	return Model{
		report: report,
	}
}
