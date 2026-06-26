package tui

import (
	"context"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (m Model) searchLocationsCmd(query string) tea.Cmd {
	query = strings.TrimSpace(query)

	return func() tea.Msg {
		if m.locationSearchLoader == nil {
			return locationSearchFailedMsg{
				err: errLocationSearchLoaderRequired,
			}
		}

		results, err := m.locationSearchLoader.SearchLocations(
			context.Background(),
			query,
		)
		if err != nil {
			return locationSearchFailedMsg{err: err}
		}

		return locationSearchLoadedMsg{results: results}
	}
}
