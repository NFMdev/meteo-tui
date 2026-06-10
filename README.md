# Meteo TUI

A terminal weather dashboard built with Go + Bubble Tea.

## Current status

v0.1.0 provides a functional interactive TUI with weather data from Open-Meteo.

## Run

```bash
go run ./cmd/meteo --city Copenhagen --country DK
```

## Controls
 - q / ctrl+c: quit
 - r: refresh
 - ↑ / k: previous day
 - ↓ / j: next day
 - ?: toggle help
