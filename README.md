# Meteo TUI

Meteo TUI is an interactive terminal weather dashboard built in Go for Linux terminal environments.

The project currently uses Open-Meteo for geocoding and forecast data.

## Current status

Version: `v0.3.0`

Meteo TUI currently supports:
- real location resolution from city and country code
- real weather forecast data from Open-Meteo
- current weather panel
- weather metrics panel
- daily forecast list
- hourly forecast for the selected day
- keyboard navigation
- XDG config path resolution
- config loading and writing

## Run

```bash
go run ./cmd/meteo --city Copenhagen --country DK
```

You can also try another city:
```bash
go run ./cmd/meteo --city Madrid --country ES
```

## Controls
| Key | Action |
| :---: | :---: |
| q / ctrl+c | quit |
| r | refresh |
| ↑ / k | previous day |
| ↓ / j | next day |
| s | open location search |
| f | open favorites screen |
| a | add current location to favorites |
| d | set current location as default |
| u / pgup | scroll up |
| d / pgdown | scroll down |
| g / home | scroll top |
| G / end | scroll bottom |
| ? | toggle help |

### Inside location search screen
| Key | Action |
| :---: | :---: |
| ↑ / ↓ | move through location results |
| Enter ⏎ | search location / load selected location |
| a | add selected location to favorites |
| d | set selected location as default |
| Esc | go back |

### Inside favorites screen
| Key | Action |
| :---: | :---: |
| ↑ / ↓ | move through favorites |
| Enter ⏎ | load weather for selected favorite |
| d | set selected favorite as default location |
| x | remove selected favorite |
| Esc | go back |

## Requirements
- Go
- Internet connection
- Terminal with enough size for the TUI

## Configuration
Meteo TUI supports a local configuration file for storing default location and favorite locations.

Default location and favorites are stored in the local config file:

`~/.config/meteo/config.json`

### Location resolution priority
Location resolution priority

Meteo resolves the final location in this order:

1. CLI flags
2. Config file
3. Built-in fallback

The built-in fallback is `Copenhagen, DK`

This means `meteo` will still run even if no config file exists.

### Location search
Users can search locations from  inside the TUI and load weather reports for selected location, they can update default location and add favorites from here.

### Initialize config
```bash
meteo --init-config --city Copenhagen --country DK
```
This creates or updates the default config file:

`~/.config/meteo/config.json`

Example config:

```json
{
  "default_city": "Copenhagen",
  "default_country": "DK"
}
```
The app will use the configured default city and country.

### Override config temporarily
CLI flags override the config file:

```bash
meteo --city Madrid --country ES
```
This uses Madrid/ES for the current run only. It does not modify the config file.

### Custom config path
You can use custom config file:
```bash
meteo --config /tmp/meteo-config.json --init-config --city Copenhagen --country DK
```
Then run with:
```bash
meteo --config /tmp/meteo-config.json
```

## Cache
Meteo TUI stores the latest successful forecast locally so the app can still be useful when the network or the weather provider are unavalilable.
By  default, forecast cache files are stored under:
`~/.cache/meteo/forecasts/`
Each location gets its own cache file.
Ex:
`~/.cache/meteo/forecasts/copenhagen_dk.json`
`~/.cache/meteo/forecasts/madrid_es.json`

### Default behavior
When running Meteo app, the program will:
1. Resolve the given location
2. Fetch data from Open-Meteo
3. Show forecast
4. Save forecast to local cache

### Offline mode
Offline mode skips network access and reads from cache only:
```bash
meteo --city Copenhagen --country DK --offline
```
[!Note] If no cache exists for the requested location, the app shows an error state.

### Custom cache directory
You can choose a custom cache directory:
```bash
meteo --city Copenhagen --country DK --cache /tmp/meteo-cache
```

## CLI flags
| Key | Action |
| :--- | :--- |
| --city | City name |
| --country | ISO 3166-1 alpha-2 country code |
| --config | Custom config file path |
| --init-config | Create or update config and exit |
| --cache | Custom forecast cache directory |
| --offline | Use cached forecast data only |
| --fail | Simulate provider failure for development/testing |

## Built With
- Go
- Bubble Tea
- Bubbles
- Lip Gloss
- Open-Meteo
