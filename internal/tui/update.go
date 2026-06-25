package tui

import (
	"context"
	"time"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadWeatherCmd(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

		switch m.mode {
		case screenModeSearchInput:
			return m.updateSearchInputKey(msg)

		default:
			return m.updateDashboardKey(msg)
		}

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case weatherLoadedMsg:
		m.loading = false
		m.err = nil
		m.report = msg.report
		m.selectedDay = 0
		m.viewport.GotoTop()
		m.rebuildViewportContent()

	case weatherFailedMsg:
		m.loading = false
		m.err = msg.err

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.SetWidth(msg.Width)
		m.configureViewport()
		m.rebuildViewportContent()
	}
	return m, nil
}

func (m Model) loadWeatherCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		report, err := m.loader.GetWeather(ctx, m.city, m.country)
		if err != nil {
			return weatherFailedMsg{err: err}
		}

		return weatherLoadedMsg{report: report}
	}
}
