package tui

import (
	"context"

	tea "charm.land/bubbletea/v2"

	"github.com/nfmdev/meteo/internal/domain"
)

func (m Model) listFavoritesCmd() tea.Cmd {
	return func() tea.Msg {
		if m.locationPreferencesLoader == nil {
			return favoritesFailedMsg{
				err: errLocationPreferencesLoaderRequired,
			}
		}

		favorites, err := m.locationPreferencesLoader.ListFavorites(context.Background())
		if err != nil {
			return favoritesFailedMsg{
				err: err,
			}
		}

		return favoritesLoadedMsg{
			favorites: favorites,
		}
	}
}

func (m Model) addFavoriteCmd(location domain.SavedLocation) tea.Cmd {
	return func() tea.Msg {
		if m.locationPreferencesLoader == nil {
			return locationPreferenceFailedMsg{
				err: errLocationPreferencesLoaderRequired,
			}
		}

		if err := m.locationPreferencesLoader.AddFavorite(
			context.Background(),
			location,
		); err != nil {
			return locationPreferenceFailedMsg{
				err: err,
			}
		}

		return locationPreferenceUpdatedMsg{
			message: "Favorite saved.",
		}
	}
}

func (m Model) removeFavoriteCmd(location domain.SavedLocation) tea.Cmd {
	return func() tea.Msg {
		if m.locationPreferencesLoader == nil {
			return locationPreferenceFailedMsg{
				err: errLocationPreferencesLoaderRequired,
			}
		}

		if err := m.locationPreferencesLoader.RemoveFavorite(
			context.Background(),
			location,
		); err != nil {
			return locationPreferenceFailedMsg{
				err: err,
			}
		}

		return locationPreferenceUpdatedMsg{
			message: "Favorite removed.",
		}
	}
}

func (m Model) setDefaultLocationCmd(location domain.SavedLocation) tea.Cmd {
	return func() tea.Msg {
		if m.locationPreferencesLoader == nil {
			return locationPreferenceFailedMsg{
				err: errLocationPreferencesLoaderRequired,
			}
		}

		if err := m.locationPreferencesLoader.SetDefaultLocation(
			context.Background(),
			location,
		); err != nil {
			return locationPreferenceFailedMsg{
				err: err,
			}
		}

		return locationPreferenceUpdatedMsg{
			message: "Default location updated.",
		}
	}
}
