package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"milesq.dev/btrbk-manage/internal/app"
	"milesq.dev/btrbk-manage/internal/ui/protect"
)

func main() {
	configPath := flag.String("c", "", "path to config file")
	project := flag.String("p", "", "project name (reads config from /etc/btrbk-manage/config.{project}.yaml)")
	flag.Parse()

	cfg, err := app.LoadConfig(*configPath, *project)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	f, _ := tea.LogToFile(filepath.Join(cfg.GetConfigPath(), "../debug.log"), "debug")
	defer f.Close()

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
