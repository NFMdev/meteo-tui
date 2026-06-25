package tui

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

func (m Model) updateSearchInputKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Cancel):
		return m.exitSearchMode(), nil

	case key.Matches(msg, m.keys.Submit):
		query := strings.TrimSpace(m.searchInput.Value())
		if query == "" {
			m.searchErr = errSearchQueryRequired
			return m, nil
		}

		m.searchErr = nil
		return m, nil
	}

	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)

	return m, cmd
}
