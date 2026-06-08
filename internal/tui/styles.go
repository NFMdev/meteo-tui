package tui

import "charm.land/lipgloss/v2"

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().Bold(true)

	subtitleStyle = lipgloss.NewStyle().Faint(true)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2)

	footerStyle = lipgloss.NewStyle().
			Faint(true)
)
