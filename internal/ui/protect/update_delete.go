package protect

import (
	tea "github.com/charmbracelet/bubbletea"
	"milesq.dev/btrbk-manage/pkg/router"
)

func (m Model) handleDeleteConfirmation(msg tea.Msg) (Model, tea.Cmd, *router.UpdateMeta) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil, router.PassThrough()
	}

	switch keyMsg.String() {
	case "y", "Y":
		m.Err = m.mng.Delete(m.selected)
		if m.Err == nil {
			m.recollect()
			if m.cursor >= len(m.backups) && m.cursor > 0 {
				m.cursor = len(m.backups) - 1
			}
		}
		m.isConfirmingDelete = false
	case "n", "N", "esc":
		m.isConfirmingDelete = false
	case "ctrl+c":
		return m, nil, router.PassThrough()
	}

	return m, nil, nil
}
