package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/nfmdev/meteo/internal/tui"
)

func main() {
	city := flag.String("city", "Copenhagen", "city name")
	country := flag.String("country", "DK", "ISO 31166-1 alpha-2 country code")

	flag.Parse()

	if *city == "" {
		fmt.Fprintln(os.Stderr, "meteo: --city cannot be empty")
		os.Exit(1)
	}

	if *country == "" {
		fmt.Fprintln(os.Stderr, "meteo: --country cannot be empty")
		os.Exit(1)
	}

	model := tui.NewModel(*city, *country)

	program := tea.NewProgram(model)

	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "meteo: failed to run TUI: %v\n", err)
		os.Exit(1)
	}
}
