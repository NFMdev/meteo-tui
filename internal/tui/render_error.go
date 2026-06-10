package tui

import (
	"fmt"
	"strings"
)

func (m Model) renderError() string {
	title := titleStyle.Render(
		fmt.Sprintf("Meteo — %s, %s", m.city, m.country),
	)

	content := errorStyle.Render(panelStyle.Render(strings.Join([]string{
		"Could not load weather data.",
		"",
		fmt.Sprintf("Error: %v", m.err),
		"",
		"Press r to retry.",
	}, "\n")))

	help := footerStyle.Render(m.help.View(m.keys))

	return appStyle.Render(strings.Join([]string{
		title,
		"",
		content,
		"",
		help,
	}, "\n"))
}
