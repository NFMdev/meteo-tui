# Changelog

## v0.2.1

### Added

- Added responsive layout for wider terminals
- Added grid panel arrangement for larger screens
- Added compact scrollable layout for smaller terminals
- Added viewport scrolling for compact mode
- Added dedicated scroll controls

### Changed

- Rearranged TUI dashboard so panels use horizontal space efficiently
- Moved from single-column dashboard layout to an adaptative layout system

## v0.2.0

### Added

- Added Open-Meteo mapper tests.
- Added Open-Meteo HTTP adapter tests using `httptest`.
- Added app service tests.
- Added terminal resize and small-screen handling.
- Added text truncation for narrow terminals.
- Added layout tests.

### Changed

- Split TUI rendering into focused files.
- Improved TUI layout.
- Improved project documentation.

## v0.1.0

### Added

- Added initial Bubble Tea TUI.
- Added current weather panel.
- Added weather metrics panel.
- Added daily forecast list.
- Added hourly forecast list for the selected day.
- Added keyboard controls.
- Added loading state.
- Added error state.
- Added Open-Meteo geocoding integration.
- Added Open-Meteo forecast integration.