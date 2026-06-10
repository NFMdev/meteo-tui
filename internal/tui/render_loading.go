package tui

import (
	"fmt"
	"strings"
)

func (m Model) renderLoading() string {
	title := titleStyle.Render(
		fmt.Sprintf("Meteo — %s, %s", m.city, m.country),
	)

	content := panelStyle.Render(strings.Join([]string{
		fmt.Sprintf("%s Loading weather data...", m.spinner.View()),
		"",
		"Fetching forecast information.",
	}, "\n"))

	help := footerStyle.Render(m.help.View(m.keys))

	return appStyle.Render(strings.Join([]string{
		title,
		"",
		content,
		"",
		help,
	}, "\n"))
}
