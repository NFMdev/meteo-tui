# Meteo TUI

Meteo TUI is an interactive terminal weather dashboard built in Go for Linux terminal environments.

The project currently uses Open-Meteo for geocoding and forecast data.

## Current status

Version: `v0.2.1`

Meteo TUI currently supports:
- real location resolution from city and country code
- real weather forecast data from Open-Meteo
- current weather panel
- weather metrics panel
- daily forecast list
- hourly forecast for the selected day
- keyboard navigation

## Run

```bash
go run ./cmd/meteo --city Copenhagen --country DK
```

You can also try another city:
```bash
go run ./cmd/meteo --city Copenhagen --country DK
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

## Built With
- Go
- Bubble Tea
- Bubbles
- Lip Gloss
- Open-Meteo
