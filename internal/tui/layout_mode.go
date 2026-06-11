package tui

type layoutMode int

const (
	layoutModeTooSmall layoutMode = iota
	layoutModeGrid
	layoutModeCompactScrollable
)

func (m Model) layoutMode() layoutMode {
	if m.width > 0 && m.height > 0 {
		if m.width < minTerminalWidth || m.height < minTerminalHeight {
			return layoutModeTooSmall
		}
	}

	if m.width >= 100 && m.height >= 28 {
		return layoutModeGrid
	}

	return layoutModeCompactScrollable
}
