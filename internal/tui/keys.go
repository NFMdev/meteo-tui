package tui

import "charm.land/bubbles/v2/key"

type KeyMap struct {
	Up      key.Binding
	Down    key.Binding
	Refresh key.Binding
	Help    key.Binding
	Quit    key.Binding

	ScrollUp     key.Binding
	ScrollDown   key.Binding
	ScrollTop    key.Binding
	ScrollBottom key.Binding

	Search key.Binding
	Cancel key.Binding
	Submit key.Binding

	AddFavorite key.Binding
	Favorites   key.Binding
}

// TODO: Separate single-column keys from multicolumn
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "previous day"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "next day"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		AddFavorite: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add favorite"),
		),
		Favorites: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "favorites"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		ScrollUp: key.NewBinding(
			key.WithKeys("pgup", "u"),
			key.WithHelp("u/pgup", "scroll up"),
		),
		ScrollDown: key.NewBinding(
			key.WithKeys("pgdown", "d"),
			key.WithHelp("d/pgdn", "scroll down"),
		),
		ScrollTop: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "top"),
		),
		ScrollBottom: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "bottom"),
		),
		Search: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "search location"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit"),
		),
	}
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Up,
		k.Down,
		k.Search,
		k.Favorites,
		k.Quit,
		k.Help,
		k.AddFavorite,
		k.Refresh,
		k.ScrollUp,
		k.ScrollDown,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Up,
			k.Down,
		},
		{
			k.ScrollUp,
			k.ScrollDown,
			k.ScrollTop,
			k.ScrollBottom,
		},
		{
			k.Refresh,
			k.Search,
		},
		{
			k.Favorites,
			k.AddFavorite,
		},
		{
			k.Help,
			k.Quit,
		},
	}
}
