package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"milesq.dev/btrbk-manage/internal/app"
	"milesq.dev/btrbk-manage/internal/ui/protect"
)

func main() {
	configPath := flag.String("c", app.DefaultConfigPath, "path to config file")
	flag.Parse()

	cfg, err := app.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	model, err := protect.InitialModel(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initializing: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())

	m, err := p.Run()

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	mm := m.(protect.Model)
	if mm.Err != nil {
		fmt.Fprintln(os.Stderr, "error:", mm.Err)
		os.Exit(1)
	}
}
