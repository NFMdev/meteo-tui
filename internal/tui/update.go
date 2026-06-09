package tui

import (
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			m.selectPreviousDay()

		case key.Matches(msg, m.keys.Down):
			m.selectNextDay()

		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
			m.help.ShowAll = m.showHelp

		case key.Matches(msg, m.keys.Refresh):
			// Temp fake refresh for Block 3.
			// In Block 4 will become an async bubbletea command.
			m.report.UpdatedAt = time.Now()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.SetWidth(msg.Width)
	}
	return m, nil
}

func (m *Model) selectPreviousDay() {
	if m.selectedDay > 0 {
		m.selectedDay--
	}
}

func (m *Model) selectNextDay() {
	if m.selectedDay < len(m.report.Daily)-1 {
		m.selectedDay++
	}
}
