package tui

import tea "charm.land/bubbletea/v2"

func (m Model) View() tea.View {
	return tea.NewView(m.render())
}

func (m Model) render() string {
	if m.isTerminalTooSmall() {
		return m.renderSmallTerminal()
	}

	switch {
	case m.loading:
		return m.renderLoading()

	case m.err != nil:
		return m.renderError()

	default:
		return m.renderDashboard()
	}
}
