package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/nfmdev/meteo/internal/app"
	forecastCache "github.com/nfmdev/meteo/internal/cache"
	meteoConfig "github.com/nfmdev/meteo/internal/config"
	"github.com/nfmdev/meteo/internal/location"
	"github.com/nfmdev/meteo/internal/tui"
	"github.com/nfmdev/meteo/internal/weather/openmeteo"
)

func main() {
	city := flag.String("city", "", "city name")
	country := flag.String("country", "", "ISO 31166-1 alpha-2 country code")
	configPath := flag.String("config", "", "custom config file path")
	cacheDir := flag.String("cache", "", "custom forecast cache directory")
	initConfig := flag.Bool("init-config", false, "create or update the config file and exit")
	offline := flag.Bool("offline", false, "use cached forecast data only")
	fail := flag.Bool("fail", false, "simulate weather loading failure")

	flag.Parse()

	resolvedOptions, err := resolveStartupOptions(startupOptions{
		city:       *city,
		country:    *country,
		configPath: *configPath,
		cacheDir:   *cacheDir,
		initConfig: *initConfig,
		offline:    *offline,
		fail:       *fail,
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

	forecastStore, err := forecastCache.NewForecastStore(resolvedOptions.cacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "meteo: failed to create forecast cache store: %v\n", err)
		os.Exit(1)
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	locationResolver := location.NewOpenMeteoResolver(httpClient)
	locationSearcher := location.NewOpenMeteoSearcher(httpClient)
	locationnSearchService := app.NewLocationSearchService(locationSearcher)

	var forecastProvider app.ForecastProvider
	if resolvedOptions.fail {
		forecastProvider = app.FailForecastProvider{}
	} else {
		forecastProvider = openmeteo.NewClient(httpClient)
	}

	weatherService := app.NewWeatherServiceWithCache(
		locationResolver,
		forecastProvider,
		forecastStore,
		app.WeatherServiceOptions{
			Offline: resolvedOptions.offline,
		},
	)

	model := tui.NewModelWithSearch(
		resolvedOptions.city,
		resolvedOptions.country,
		weatherService,
		locationnSearchService,
	)

	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "meteo: failed to run TUI: %v\n", err)
		os.Exit(1)
	}
}
