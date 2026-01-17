package protect

import (
	tea "github.com/charmbracelet/bubbletea"
	"milesq.dev/btrbk-manage/pkg/multiselect"
	"milesq.dev/btrbk-manage/pkg/router"
)

func (m Model) handleRestoreSelector(msg tea.Msg) (Model, tea.Cmd, *router.UpdateMeta) {
	exitMsg, isExitMsg := msg.(multiselect.ExitMsg)

	if isExitMsg {
		m.isChoosingSubvolumesForRestore = false

		if exitMsg.Reason == multiselect.UserSaved && len(exitMsg.Selected) > 0 {
			m.Err = m.mng.Restore(m.selected, exitMsg.Selected)
		}
		return m, nil, nil
	}

	var cmd tea.Cmd
	var meta *router.UpdateMeta
	m.restoreSelector, cmd, meta = m.restoreSelector.Update(msg)
	return m, cmd, meta
}
