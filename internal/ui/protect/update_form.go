package protect

import (
	tea "github.com/charmbracelet/bubbletea"
	"milesq.dev/btrbk-manage/pkg/form"
	"milesq.dev/btrbk-manage/pkg/router"
)

func (m Model) handleForm(msg tea.Msg) (Model, tea.Cmd, *router.UpdateMeta) {
	exitMsg, isExitMsg := msg.(form.ExitMsg)
	if isExitMsg {
		m.isEdit = false

		if exitMsg.Reason == form.UserSaved {
			m.handleSave(exitMsg.Values)
		}
		return m, nil, nil
	}

	var formCmd tea.Cmd
	var meta *router.UpdateMeta
	m.form, formCmd, meta = m.form.Update(msg)
	return m, formCmd, meta
}

func (m *Model) handleSave(values []string) {
	note, err := m.getProtectionNote(values)
	if err != nil {
		m.Err = err
		return
	}
	m.Err = m.mng.Protect(m.selected.Timestamp, note)
	if m.Err != nil {
		return
	}
	m.recollect()
}
