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

func (m Model) isTerminalTooSmall() bool {
	if m.width <= 0 || m.height <= 0 {
		return false
	}

	return m.width < minTerminalWidth || m.height < minTerminalHeight
}
