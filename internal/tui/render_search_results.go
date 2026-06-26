package tui

import "fmt"

func (m Model) renderSearchResults() string {
	lines := []string{
		"Search Results",
		"",
	}

	if m.searching {
		lines = append(lines, m.spinner.View()+" Searching locations...")
		lines = append(lines, "", "Press Esc to cancel")
		return panelStyle.
			Width(m.panelWidth()).
			Render(joinTruncatedLines(lines, m.innerPanelWidth()))
	}

	if m.searchErr != nil {
		lines = append(lines, "Error: "+m.searchErr.Error())
		lines = append(lines, "", "Press Esc to return")
		return panelStyle.
			Width(m.panelWidth()).
			Render(joinTruncatedLines(lines, m.innerPanelWidth()))
	}

	if len(m.searchResults) == 0 {
		lines = append(lines, "No results found.")
		lines = append(lines, "", "Press Esc to search again")
		return panelStyle.
			Width(m.panelWidth()).
			Render(joinTruncatedLines(lines, m.innerPanelWidth()))
	}

	for index, result := range m.searchResults {
		cursor := " "
		if index == m.selectedSearchResult {
			cursor = ">"
		}

		lines = append(
			lines,
			fmt.Sprintf("%s %s", cursor, locationSearchResultLabel(result)),
		)
	}

	lines = append(lines, "", "↑/↓ select result • Enter ⏎ select • Esc back")
	return panelStyle.
		Width(m.panelWidth()).
		Render(joinTruncatedLines(lines, m.innerPanelWidth()))
}
