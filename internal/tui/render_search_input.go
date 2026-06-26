package tui

func (m Model) renderSearchInput() string {
	lines := []string{
		"Search location",
		"",
		m.searchInput.View(),
		"",
		"Type a city name and press Enter",
		"Press Esc to cancel",
	}

	if m.searchErr != nil {
		lines = append(lines, "", "Error: "+m.searchErr.Error())
	}

	return panelStyle.Width(m.panelWidth()).Render(joinTruncatedLines(lines, m.innerPanelWidth()))
}
