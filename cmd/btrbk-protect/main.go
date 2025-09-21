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

	if mm, ok := m.(protect.Model); ok && mm.Selected != nil {
		for _, it := range mm.Selected.Items {
			fmt.Println(it.Path)
		}
	}
}
