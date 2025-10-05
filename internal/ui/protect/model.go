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
		form: form.New(inputs, form.FormStyles{
			BlurredButton: blurredButton,
			FocusedButton: focusedButton,
			BlurStyle:     lipgloss.NewStyle(),
			FocuseStyle:   focusedStyle,
		}),
	}
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var formCmd tea.Cmd

	if m.SelectedForEdit != nil {
		var meta *components.UpdateMeta
		m.form, formCmd, meta = m.form.Update(msg)

		if meta != nil && meta.Finish {
			m.SelectedForEdit = nil
		}

		if meta == nil || !meta.PassThrough {
			return m, formCmd
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Batch(formCmd, tea.Quit)
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
		case "enter":
			if m.Err == nil && len(m.Groups) > 0 {
				m.SelectedForEdit = &m.Groups[m.Cursor]
			}
		}
	}
	return m, formCmd
}
