package protect

import (
	"fmt"
	"strings"

	"milesq.dev/btrbk-manage/internal/utils"
)

func (m Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf("Error: %v\n\nDir: %s\nPress q to quit.\n", m.Err, m.Dir)
	}
	var b strings.Builder
	title := fmt.Sprintf("Btrbk backups in %s  —  %d backups\n", m.Dir, len(m.Backups))
	b.WriteString(title)
	b.WriteString(strings.Repeat("─", utils.MinMax(10, len(title), 80)))
	b.WriteString("\n\n")

	if m.IsConfirmingDelete {
		m.viewDeleteConfirmation(&b)
	} else if m.IsEdit {
		b.WriteString(m.form.View())
	} else {
		m.viewList(&b)
		writeHelpMessage(&b)
	}

	return b.String()
}
