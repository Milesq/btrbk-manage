package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"milesq.dev/btrbk-manage/internal/ui/protect"
)

func main() {
	p := tea.NewProgram(protect.InitialModel("./mnt/@snaps"), tea.WithAltScreen())

	m, err := p.Run()

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	mm := m.(protect.Model)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", mm.Err)
		os.Exit(1)
	}
}
