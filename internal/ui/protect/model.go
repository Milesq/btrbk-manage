package protect

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"milesq.dev/btrbk-manage/internal/snaps"
	"milesq.dev/btrbk-manage/pkg/components"
	"milesq.dev/btrbk-manage/pkg/form"
)

type Model struct {
	// Core
	mng            snaps.BackupManager
	Dir            string
	Groups         []snaps.Group
	TotalSnapshots int

	// General State
	Err    error
	Cursor int

	// Modes flags
	ListProtectedOnly bool
	TrashMode         bool
	SelectedForEdit   *snaps.Group

	// SubComponents
	form form.Model
}

func InitialModel(dir string) Model {
	backupManager := snaps.GetManagerForDirectory(dir)
	backups, err := backupManager.Collect()

	inputs := getProtectionNoteInputs()

	return Model{
		Groups:         backups.Groups,
		Err:            err,
		Dir:            dir,
		TotalSnapshots: backups.TotalCount,
		mng:            backupManager,
		form: form.New(inputs, form.NewFormProps().WithStyles(form.FormStyles{
			BlurredButton: blurredButton,
			FocusedButton: focusedButton,
			BlurStyle:     lipgloss.NewStyle(),
			FocuseStyle:   focusedStyle,
		})),
	}
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	if m.SelectedForEdit != nil {
		var formCmd tea.Cmd
		var meta *components.UpdateMeta
		m.form, formCmd, meta = m.form.Update(msg)

		if meta == nil || !meta.PassThrough {
			return m, formCmd
		}
		cmds = append(cmds, formCmd)
	}

	switch msg := msg.(type) {
	case form.ExitMsg:
		if msg.Reason == form.UserSaved {
			note, err := snaps.GetProtectionNote(msg.Values)
			if err != nil {
				m.Err = err
				return m, tea.Batch(cmds...)
			}

			m.Err = m.mng.Protect(m.SelectedForEdit.Timestamp, note)

			m.recollect()
		}
		m.SelectedForEdit = nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			cmds = append(cmds, tea.Quit)
		case "up", "k":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = max(0, len(m.Groups)-1)
			}
		case "down", "j":
			m.Cursor++
			if m.Cursor >= len(m.Groups) {
				m.Cursor = 0
			}
		case " ":
			group := &m.Groups[m.Cursor]
			if group.IsProtected {
				m.Err = m.mng.FreePersistance(group.Timestamp)
				m.recollect()
			} else {
				m.SelectedForEdit = group
			}
		case "enter":
			if m.Err == nil && len(m.Groups) > 0 {
				group := &m.Groups[m.Cursor]
				if group.IsProtected {
					m.SelectedForEdit = group
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) recollect() {
	m.mng.ClearCache()
	backups, err := m.mng.Collect()
	m.Err = err
	m.TotalSnapshots = backups.TotalCount
	m.Groups = backups.Groups
}
