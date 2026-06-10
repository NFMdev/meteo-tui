package tui

import (
	"fmt"
	"strings"
)

func (m Model) renderSmallTerminal() string {
	title := titleStyle.Render("Meteo")

	content := panelStyle.
		Width(maxInt(30, m.contentWidth())).
		Render(strings.Join([]string{
			"Terminal too small.",
			"",
			fmt.Sprintf("Current: %dx%d", m.width, m.height),
			fmt.Sprintf("Minimum: %dx%d", minTerminalWidth, minTerminalHeight),
			"",
			"Please resize your terminal.",
			"",
			"q quit • ctrl+c quit",
		}, "\n"))

	return appStyle.Render(strings.Join([]string{
		title,
		"",
		content,
	}, "\n"))
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
