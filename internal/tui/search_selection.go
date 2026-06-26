package tui

import "github.com/nfmdev/meteo/internal/domain"

func (m Model) selectedLocationSearchResult() (domain.LocationSearchResult, bool) {
	if len(m.searchResults) == 0 ||
		m.selectedSearchResult < 0 ||
		m.selectedSearchResult >= len(m.searchResults) {
		return domain.LocationSearchResult{}, false
	}

	return m.searchResults[m.selectedSearchResult], true
}
