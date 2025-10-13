package protect

import (
	tea "github.com/charmbracelet/bubbletea"
	"milesq.dev/btrbk-manage/pkg/router"
)

func (m Model) handleList(msg tea.Msg) (Model, tea.Cmd, *router.UpdateMeta) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit, nil
		case "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = max(0, len(m.Backups)-1)
			}
		case "down", "j":
			m.Cursor++
			if m.Cursor >= len(m.Backups) {
				m.Cursor = 0
			}
		case " ":
			if len(m.Backups) > 0 {
				backup := m.Backups[m.Cursor]
				if backup.IsProtected {
					m.Err = m.mng.FreePersistance(backup.Timestamp)
					m.recollect()
				} else {
					m.selected = backup
					m.IsEdit = true
					m.populateFormWithNote(backup.ProtectionNote)
				}
			}
		case "enter":
			if m.Err == nil && len(m.Backups) > 0 {
				backup := m.Backups[m.Cursor]
				if backup.IsProtected {
					m.selected = backup
					m.IsEdit = true
					m.populateFormWithNote(backup.ProtectionNote)
				}
			}
		case "d":
			if m.Err == nil && len(m.Backups) > 0 {
				m.selected = m.Backups[m.Cursor]
				m.IsConfirmingDelete = true
			}
		}
	}

	return m, nil, nil
}
