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
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			m.selectPreviousDay()

		case key.Matches(msg, m.keys.Down):
			m.selectNextDay()

		case key.Matches(msg, m.keys.ScrollUp):
			if m.layoutMode() == layoutModeCompactScrollable {
				m.viewport.ScrollUp(3)
			}

		case key.Matches(msg, m.keys.ScrollDown):
			if m.layoutMode() == layoutModeCompactScrollable {
				m.viewport.ScrollDown(3)
			}

		case key.Matches(msg, m.keys.ScrollTop):
			if m.layoutMode() == layoutModeCompactScrollable {
				m.viewport.GotoTop()
			}

		case key.Matches(msg, m.keys.ScrollBottom):
			if m.layoutMode() == layoutModeCompactScrollable {
				m.viewport.GotoBottom()
			}

		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
			m.help.ShowAll = m.showHelp

		case key.Matches(msg, m.keys.Refresh):
			m.loading = true
			m.err = nil
			m.selectedDay = 0

			return m, tea.Batch(
				m.spinner.Tick,
				m.loadWeatherCmd(),
			)
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

func (m *Model) selectPreviousDay() {
	if m.loading || m.err != nil {
		return
	}

	if m.selectedDay > 0 {
		m.selectedDay--
	}
}

func (m *Model) selectNextDay() {
	if m.loading || m.err != nil {
		return
	}

	if m.selectedDay < len(m.report.Daily)-1 {
		m.selectedDay++
	}
}
