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
| u / pgup | scroll up |
| d / pgdown | scroll down |
| g / home | scroll top |
| G / end | scroll bottom |
| ? | toggle help |

## Requirements
- Go
- Internet connection
- Terminal with enough size for the TUI

## Configuration
Meteo TUI supports a local configuration file for storing a default location.

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

### Location resolution priority
Location resolution priority

Meteo resolves the final location in this order:

1. CLI flags
2. Config file
3. Built-in fallback

The built-in fallback is:
Copenhagen, DK

This means:
```bash
meteo
```
will still work even if no config file exists.

### CLI flags
| Key | Action |
| :--- | :--- |
| --city | City name |
| --country | ISO 3166-1 alpha-2 country code |
| --config | Custom config file path |
| --init-config | Create or update config and exit |
| --fail | Simulate weather loading failure for development/testing |

## Built With
- Go
- Bubble Tea
- Bubbles
- Lip Gloss
- Open-Meteo
