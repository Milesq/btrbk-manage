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
		// m.Err = m.mng.Delete(m.selected)
		if m.Err == nil {
			m.recollect()
			if m.Cursor >= len(m.Backups) && m.Cursor > 0 {
				m.Cursor = len(m.Backups) - 1
			}
		}
		m.IsConfirmingDelete = false
	case "n", "N", "esc":
		m.IsConfirmingDelete = false
	case "ctrl+c":
		return m, nil, router.PassThrough()
	}

	return m, nil, nil
}
