package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type snapshot struct {
	Path      string
	BaseName  string
	Timestamp string // e.g., 20250810T102000 (no trailing Z)
}

type group struct {
	Timestamp string
	Items     []snapshot
}

type model struct {
	groups         []group
	cursor         int
	width, height  int
	ready          bool
	err            error
	selected       *group
	dir            string
	totalSnapshots int
}

var tsRe = regexp.MustCompile(`(?i)(\d{8}T\d{6})Z?$`)

// detectSnapshot returns (timestamp, ok).
func detectSnapshot(base string) (string, bool) {
	// Expect something like "<anything>.<YYYYMMDDTHHMMSS>[Z]"
	// Find last dot, take the suffix after it.
	dot := strings.LastIndexByte(base, '.')
	if dot < 0 || dot == len(base)-1 {
		return "", false
	}
	suffix := base[dot+1:]
	m := tsRe.FindStringSubmatch(suffix)
	if m == nil {
		return "", false
	}
	return m[1], true
}

func collect(dir string) ([]group, int, error) {
	gmap := make(map[string][]snapshot)
	var total int

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Skip unreadable entries but continue.
			return nil
		}
		// Only consider directories (typical for btrfs snapshots), but allow files too in case of send streams, etc.
		base := filepath.Base(path)
		ts, ok := detectSnapshot(base)
		if !ok {
			return nil
		}
		// Optional: avoid descending into nested snapshots too deep—WalkDir already handles recursion.
		info, e := d.Info()
		if e == nil && info.Mode().IsDir() || (e == nil && info.Mode().IsRegular()) {
			gmap[ts] = append(gmap[ts], snapshot{
				Path:      path,
				BaseName:  base,
				Timestamp: ts,
			})
			total++
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	if len(gmap) == 0 {
		return nil, 0, errors.New("no snapshots matching pattern *.YYYYMMDDTHHMMSS[Z] found")
	}

	// Build groups slice and sort members by name.
	groups := make([]group, 0, len(gmap))
	for ts, items := range gmap {
		sort.Slice(items, func(i, j int) bool { return items[i].BaseName < items[j].BaseName })
		groups = append(groups, group{Timestamp: ts, Items: items})
	}
	// Sort groups by timestamp desc (string compare works for this format).
	sort.Slice(groups, func(i, j int) bool { return groups[i].Timestamp > groups[j].Timestamp })

	return groups, total, nil
}

// ----- Bubble Tea -----

func initialModel(dir string) model {
	groups, total, err := collect(dir)
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
			return m, nil
		case "down", "j":
			if len(m.groups) > 0 && m.cursor < len(m.groups)-1 {
				m.cursor++
			}
			return m, nil
		case "enter":
			if m.err == nil && len(m.groups) > 0 {
				g := m.groups[m.cursor]
				m.selected = &g
				return m, tea.Quit
			}
			return m, nil
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

// ----- main -----

func main() {
	dir := flag.String("dir", ".", "Directory to scan recursively for snapshots")
	flag.Parse()

	p := tea.NewProgram(initialModel(*dir), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	// After TUI exits, if a group was selected, print its members (one per line).
	if mm, ok := m.(model); ok && mm.selected != nil {
		for _, it := range mm.selected.Items {
			fmt.Println(it.Path)
		}
	}
}

// ----- utils -----
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
