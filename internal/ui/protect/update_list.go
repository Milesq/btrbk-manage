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
			m.cursor--
			if m.cursor < 0 {
				m.cursor = max(0, len(m.backups)-1)
			}
		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.backups) {
				m.cursor = 0
			}
		case " ":
			if len(m.backups) > 0 {
				backup := m.backups[m.cursor]
				if backup.IsProtected {
					m.Err = m.mng.FreePersistance(backup.Timestamp)
					m.recollect()
				} else {
					m.selected = backup
					m.isEdit = true
					m.populateFormWithNote(backup.ProtectionNote)
				}
			}
		case "enter":
			if m.Err == nil && len(m.backups) > 0 {
				backup := m.backups[m.cursor]
				if backup.IsProtected {
					m.selected = backup
					m.isEdit = true
					m.populateFormWithNote(backup.ProtectionNote)
				}
			}
		case "d":
			if m.Err == nil && len(m.backups) > 0 {
				m.selected = m.backups[m.cursor]
				m.isConfirmingDelete = true
			}
		case "t":
			m.trashMode = !m.trashMode
			if m.trashMode {
				m.listProtectedOnly = false
			}
			m.cursor = 0
			m.recollect()
		case "m":
			m.listProtectedOnly = !m.listProtectedOnly
			if m.listProtectedOnly {
				m.trashMode = false
			}
			m.cursor = 0
			m.recollect()
		case "D":
			if m.trashMode && len(m.backups) > 0 {
				for _, backup := range m.backups {
					if err := m.mng.RemoveFromTrash(backup); err != nil {
						m.Err = err
						break
					}
				}
				m.trashMode = false
				m.recollect()
			}
		case "r":
			if m.Err == nil && len(m.backups) > 0 {
				m.selected = m.backups[m.cursor]
				m.isChoosingSubvolumesForRestore = true
				m.restoreSelector = m.createRestoreSelector()
			}
		}
	}

	return m, nil, nil
}
