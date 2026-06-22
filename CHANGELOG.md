# Changelog

## v0.4.0

### Added
- Added local cache
- Added XDG default cache path resolution
- Added JSON cache files
- Added key generation from city and country, giving a cache file for each city-country
- Added cache store for reading and writing
- Added `--cache` CLI flag for custom cache directories
- Added `__offline` CLI flag to skp network accesss and use cache only
- Added cached-data indicator in the TUI header
- Added tests for cache keys, paths, store, source metadata, app fallback behavior, and source labels

### Changed
- The weather service can now save successful forecasts to cache
- The weather service can fallback to cache if location or forecast loading fails
- `--fail` development flag now simulates provider failure so cache fallback can be tested
- TUI header displays whether data is new or cached
- Changed grid content height to fit header 

## v0.3.0

### Added
- Added config model for default city and country.
- Added XDG default config path resolution.
- Added support for custom config paths with `--config`.
- Added config loading from JSON.
- Added config writing to JSON.
- Added `--init-config` to create or update local config.
- Added startup resolution priority: CLI flags, config file, fallback.
- Added validation for default city and country.
- Added country code normalization to uppercase.
- Added tests for config path resolution.
- Added tests for config loading and writing.
- Added tests for startup option resolution.

### Changed
- `--city` and `--country` now default to empty values internally so the app can distinguish explicit CLI overrides from config values.
- meteo can now run without `--city` and `--country`.
- Startup logic now resolves the final location before creating the TUI model.

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