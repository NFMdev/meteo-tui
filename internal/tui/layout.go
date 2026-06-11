package tui

const (
	minTerminalWidth  = 80
	minTerminalHeight = 24

	defaultTerminalWidth = 100

	appHorizontalPadding = 4
	panelHorizontalFrame = 6
)

func (m Model) terminalWidth() int {
	if m.width <= 0 {
		return defaultTerminalWidth
	}

	return m.width
}

func (m Model) contentWidth() int {
	width := m.terminalWidth() - appHorizontalPadding
	if width < 1 {
		return 1
	}

	return width
}

func (m Model) panelWidth() int {
	width := m.contentWidth()
	if width < 1 {
		return 1
	}

	return width
}

func (m Model) innerPanelWidth() int {
	width := m.panelWidth() - panelHorizontalFrame
	if width < 1 {
		return 1
	}

	return width
}

func (m Model) compactViewportHeight() int {
	if m.height <= 0 {
		return 20
	}

	// Approx:
	// app vertical padding + header + blank lines + help/footer.
	height := m.height - 7

	if height < 5 {
		return 5
	}

	return height
}
