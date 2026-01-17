package protect

import (
	"fmt"
	"strings"

	"milesq.dev/btrbk-manage/internal/utils"
)

func (m Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf("Error: %v\n\nDir: %s\nPress q to quit.\n", m.Err, m.dir)
	}
	var b strings.Builder
	title := fmt.Sprintf("Btrbk backups in %s  —  %d backups\n", m.dir, len(m.backups))
	b.WriteString(title)
	b.WriteString(strings.Repeat("─", utils.MinMax(10, len(title), 80)))
	b.WriteString("\n\n")

	if m.isConfirmingDelete {
		m.viewDeleteConfirmation(&b)
	} else if m.isEdit {
		b.WriteString(m.form.View())
	} else if m.isChoosingSubvolumesForRestore {
		b.WriteString(m.restoreSelector.View())
	} else {
		m.viewList(&b)
		m.writeHelpMessage(&b)
	}

	return b.String()
}
