package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"milesq.dev/btrbk-manage/internal/snaps"
)

type model struct {
	groups         []snaps.Group
	cursor         int
	width, height  int
	ready          bool
	err            error
	selected       *snaps.Group
	dir            string
	totalSnapshots int
}

func initialModel(dir string) model {
	groups, total, err := snaps.Collect(dir)
	return model{
		groups:         groups,
		cursor:         0,
		err:            err,
		dir:            dir,
		totalSnapshots: total,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if len(m.groups) > 0 && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if len(m.groups) > 0 && m.cursor < len(m.groups)-1 {
				m.cursor++
			}
		case "enter":
			if m.err == nil && len(m.groups) > 0 {
				g := m.groups[m.cursor]
				m.selected = &g
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if !m.ready {
		return "Loading…\n"
	}
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nDir: %s\nPress q to quit.\n", m.err, m.dir)
	}
	var b strings.Builder
	title := fmt.Sprintf("Btrbk snapshots in %s  —  %d groups, %d items\n", m.dir, len(m.groups), m.totalSnapshots)
	b.WriteString(title)
	b.WriteString(strings.Repeat("─", max(10, min(len(title), 80))))
	b.WriteString("\n\n")

	if len(m.groups) == 0 {
		b.WriteString("No snapshot groups found.\n")
		return b.String()
	}

	// Display groups with counts; highlight cursor
	for i, g := range m.groups {
		line := fmt.Sprintf("  %s  (%d)", g.Timestamp, len(g.Items))
		if i == m.cursor {
			line = "> " + line[2:]
		}
		b.WriteString(line + "\n")
	}
	b.WriteString("\n↑/k, ↓/j to move • Enter to select • q/Esc to quit\n")
	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel("./mnt/@snaps"), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if mm, ok := m.(model); ok && mm.selected != nil {
		for _, it := range mm.selected.Items {
			fmt.Println(it.Path)
		}
	}
}
