package tui

import "fmt"

func (m Model) renderFavorites() string {
	lines := []string{
		"Favorites",
		"",
	}

	if m.favoritesLoading {
		lines = append(lines, m.spinner.View()+" Loading favorites...")
		lines = append(lines, "", "Press Esc to return")
		return panelStyle.
			Width(m.panelWidth()).
			Render(joinTruncatedLines(lines, m.innerPanelWidth()))
	}

	if m.favoritesErr != nil {
		lines = append(lines, "Error: "+m.favoritesErr.Error())
		lines = append(lines, "", "Press Esc to return")
		return panelStyle.
			Width(m.panelWidth()).
			Render(joinTruncatedLines(lines, m.innerPanelWidth()))
	}

	if len(m.favorites) == 0 {
		lines = append(lines, "No favorites saved yet")
		lines = append(lines, "", "Press `a` from the dashboard screen to save the current location")
		lines = append(lines, "", "Press Esc to return")
		return panelStyle.
			Width(m.panelWidth()).
			Render(joinTruncatedLines(lines, m.innerPanelWidth()))
	}

	for index, favorite := range m.favorites {
		cursor := " "
		if index == m.selectedFavorite {
			cursor = ">"
		}
		lines = append(
			lines,
			fmt.Sprintf("%s %s", cursor, savedLocationLabel(favorite)),
		)
	}

	lines = append(
		lines,
		"",
		"↑/↓ select favorite • Enter ⏎ load weather • d set default • x remove • Esc back",
	)

	return panelStyle.
		Width(m.panelWidth()).
		Render(joinTruncatedLines(lines, m.innerPanelWidth()))
}
