package tui

type dashboardLayout struct {
	contentWidth  int
	contentHeight int

	leftWidth  int
	rightWidth int
	gap        int

	currentHeight int
	metricsHeight int
	dailyHeight   int
	hourlyHeight  int
}

func (m Model) dashboardLayout() dashboardLayout {
	gap := 2

	contentWidth := m.contentWidth()
	contentHeight := m.dashboardContentHeight()

	leftWidth := contentWidth * 40 / 100
	rightWidth := contentWidth - leftWidth - gap

	if leftWidth < 30 {
		leftWidth = 30
	}

	if rightWidth < 40 {
		rightWidth = 40
	}

	currentHeight := 7
	metricsHeight := 9

	dailyHeight := contentHeight - currentHeight - metricsHeight
	if dailyHeight < 8 {
		dailyHeight = 8
	}

	hourlyHeight := contentHeight

	return dashboardLayout{
		contentWidth:  contentWidth,
		contentHeight: contentHeight,

		leftWidth:  leftWidth,
		rightWidth: rightWidth,
		gap:        gap,

		currentHeight: currentHeight,
		metricsHeight: metricsHeight,
		dailyHeight:   dailyHeight,
		hourlyHeight:  hourlyHeight,
	}
}

func (m Model) dashboardContentHeight() int {
	if m.height <= 0 {
		return 26
	}

	height := m.height - 7
	if height < 10 {
		return 10
	}

	return height
}
