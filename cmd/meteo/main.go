package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/nfmdev/meteo/internal/app"
	meteoConfig "github.com/nfmdev/meteo/internal/config"
	"github.com/nfmdev/meteo/internal/location"
	"github.com/nfmdev/meteo/internal/tui"
	"github.com/nfmdev/meteo/internal/weather/openmeteo"
)

func main() {
	city := flag.String("city", "", "city name")
	country := flag.String("country", "", "ISO 31166-1 alpha-2 country code")
	configPath := flag.String("config", "", "custom config file path")
	initConfig := flag.Bool("init-config", false, "create or update the config file and exit")
	fail := flag.Bool("fail", false, "simulate weather loading failure")
	fake := flag.Bool("fake", false, "simulate fake weather service")

	flag.Parse()

	resolvedOptions, err := resolveStartupOptions(startupOptions{
		city:       *city,
		country:    *country,
		configPath: *configPath,
		initConfig: *initConfig,
		fail:       *fail,
		fake:       *fake,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "meteo: %v\n", err)
		os.Exit(1)
	}

	if resolvedOptions.initConfig {
		err := meteoConfig.WriteConfig(resolvedOptions.configPath, meteoConfig.AppConfig{
			DefaultCity:    resolvedOptions.city,
			DefaultCountry: resolvedOptions.country,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "meteo: failed to write config: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(
			os.Stdout,
			"meteo: config written to %s\nDefault location: %s, %s\n",
			resolvedOptions.configPath,
			resolvedOptions.city,
			resolvedOptions.country,
		)

		return
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	locationResolver := location.NewOpenMeteoResolver(httpClient)
	forecastClient := openmeteo.NewClient(httpClient)

	var weatherService app.WeatherService
	if *fail {
		weatherService = app.NewFailingWeatherService(locationResolver)
	} else if *fake {
		weatherService = app.NewFakeWeatherService(locationResolver)
	} else {
		weatherService = app.NewWeatherService(locationResolver, forecastClient)
	}

	model := tui.NewModel(resolvedOptions.city, resolvedOptions.country, weatherService)

	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "meteo: failed to run TUI: %v\n", err)
		os.Exit(1)
	}
}
