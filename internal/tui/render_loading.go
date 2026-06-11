package tui

import (
	"fmt"
	"strings"
)

func (m Model) renderLoading() string {
	title := titleStyle.Render(
		truncateText(fmt.Sprintf("Meteo — %s, %s", m.city, m.country), m.contentWidth()),
	)

	lines := []string{
		fmt.Sprintf("%s Loading weather data...", m.spinner.View()),
		"",
		"Fetching forecast information.",
	}

	content := panelStyle.
		Width(m.panelWidth()).
		Render(joinTruncatedLines(lines, m.innerPanelWidth()))

	help := footerStyle.Render(m.help.View(m.keys))

	return appStyle.Render(strings.Join([]string{
		title,
		"",
		content,
		"",
		help,
	}, "\n"))
}
