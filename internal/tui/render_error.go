package tui

import (
	"fmt"
	"strings"
)

func (m Model) renderError() string {
	title := titleStyle.Render(
		truncateText(fmt.Sprintf("Meteo — %s, %s", m.city, m.country), m.contentWidth()),
	)

	lines := []string{
		"Could not load weather data.",
		"",
		fmt.Sprintf("Error: %v", m.err),
		"",
		"Press r to retry.",
	}

	content := errorStyle.
		Width(m.panelWidth()).
		Render(panelStyle.Render(joinTruncatedLines(lines, m.innerPanelWidth())))

	help := footerStyle.Render(truncateText(m.help.View(m.keys), m.contentWidth()))

	return appStyle.Render(strings.Join([]string{
		title,
		"",
		content,
		"",
		help,
	}, "\n"))
}
