package protect

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"milesq.dev/btrbk-manage/internal/snaps"
	"milesq.dev/btrbk-manage/internal/utils"
)

type Model struct {
	Dir string
	mng snaps.BackupManager

	Err    error
	Cursor int

	Selected       *snaps.Group
	Groups         []snaps.Group
	TotalSnapshots int
}

func InitialModel(dir string) Model {
	backupManager := snaps.GetManagerForDirectory(dir)
	backups, err := backupManager.Collect()

	return Model{
		Groups:         backups.Groups,
		Cursor:         0,
		Err:            err,
		Dir:            dir,
		TotalSnapshots: backups.TotalCount,
		mng:            backupManager,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if len(m.Groups) > 0 && m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if len(m.Groups) > 0 && m.Cursor < len(m.Groups)-1 {
				m.Cursor++
			}
		case "enter":
			if m.Err == nil && len(m.Groups) > 0 {
				g := m.Groups[m.Cursor]
				m.Selected = &g
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf("Error: %v\n\nDir: %s\nPress q to quit.\n", m.Err, m.Dir)
	}
	var b strings.Builder
	title := fmt.Sprintf("Btrbk backups in %s  —  %d backups, %d snapshots\n", m.Dir, len(m.Groups), m.TotalSnapshots)
	b.WriteString(title)
	b.WriteString(strings.Repeat("─", utils.MinMax(10, len(title), 80)))
	b.WriteString("\n\n")

	if len(m.Groups) == 0 {
		b.WriteString("No snapshot groups found.\n")
		return b.String()
	}

	for i, g := range m.Groups {
		line := fmt.Sprintf("  %s", utils.PrettifyDate(g.Timestamp))
		if i == m.Cursor {
			line = "> " + line[2:]
		}
		b.WriteString(line + "\n")
	}
	b.WriteString("\n↑/↓ to move • Enter to select • q to quit\n")
	return b.String()
}
