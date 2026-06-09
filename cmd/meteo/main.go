package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/nfmdev/meteo/internal/app"
	"github.com/nfmdev/meteo/internal/location"
	"github.com/nfmdev/meteo/internal/tui"
)

func main() {
	city := flag.String("city", "Copenhagen", "city name")
	country := flag.String("country", "DK", "ISO 31166-1 alpha-2 country code")
	fail := flag.Bool("fail", false, "simulate weather loading failure")

	flag.Parse()

	if *city == "" {
		fmt.Fprintln(os.Stderr, "meteo: --city cannot be empty")
		os.Exit(1)
	}

	if *country == "" {
		fmt.Fprintln(os.Stderr, "meteo: --country cannot be empty")
		os.Exit(1)
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	locationResolver := location.NewOpenMeteoResolver(httpClient)

	var weatherService app.WeatherService
	if *fail {
		weatherService = app.NewFailingWeatherService(locationResolver)
	} else {
		weatherService = app.NewFakeWeatherService(locationResolver)
	}

	model := tui.NewModel(*city, *country, weatherService)

	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "meteo: failed to run TUI: %v\n", err)
		os.Exit(1)
	}
}
