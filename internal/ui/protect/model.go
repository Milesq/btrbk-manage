package protect

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"milesq.dev/btrbk-manage/internal/snaps"
	"milesq.dev/btrbk-manage/pkg/components"
	"milesq.dev/btrbk-manage/pkg/form"
)

type Model struct {
	// Core
	mng     snaps.BackupManager
	Dir     string
	Backups []snaps.Backup

	// General State
	Err             error
	Cursor          int
	SelectedForEdit snaps.Backup

	// Modes flags
	ListProtectedOnly bool
	TrashMode         bool
	IsEdit            bool

	// SubComponents
	form form.Model
}

func InitialModel(dir string) Model {
	backupManager := snaps.GetManagerForDirectory(dir)
	info, err := backupManager.Collect()

	inputs := getProtectionNoteInputs()

	return Model{
		Backups: info.Backups,
		Err:     err,
		Dir:     dir,
		mng:     backupManager,
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
	if m.IsEdit {
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
			note, err := m.getProtectionNote(msg.Values)
			if err != nil {
				m.Err = err
				return m, tea.Batch(cmds...)
			}

			m.Err = m.mng.Protect(m.SelectedForEdit.Timestamp, note)

			m.recollect()
		}
		m.IsEdit = false
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			cmds = append(cmds, tea.Quit)
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
			backup := m.Backups[m.Cursor]
			if backup.IsProtected {
				m.Err = m.mng.FreePersistance(backup.Timestamp)
				m.recollect()
			} else {
				m.SelectedForEdit = backup
				m.IsEdit = true
				m.populateFormWithNote(backup.ProtectionNote)
			}
		case "enter":
			if m.Err == nil && len(m.Backups) > 0 {
				backup := m.Backups[m.Cursor]
				if backup.IsProtected {
					m.SelectedForEdit = backup
					m.IsEdit = true
					m.populateFormWithNote(backup.ProtectionNote)
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
	m.Backups = backups.Backups
}

const building_params_n int = 3

func (m *Model) populateFormWithNote(note snaps.ProtectionNote) {
	if len(m.form.Inputs) < 3 {
		fmt.Println("skip")
		return
	}

	m.form.Inputs[0].SetValue(note.Note)
	m.form.Inputs[1].SetValue(note.Reason)
	if len(note.Tags) > 0 {
		m.form.Inputs[2].SetValue(strings.Join(note.Tags, ", "))
	}
}

func (m *Model) getProtectionNote(values []string) (note snaps.ProtectionNote, err error) {
	if len(values) != building_params_n {
		return snaps.ProtectionNote{}, fmt.Errorf("building protection note failed. Expected %d params, got %d", building_params_n, len(values))
	}
	note.Note = values[0]
	note.Reason = values[1]
	note.Tags = strings.Split(values[2], ",")
	for j := range note.Tags {
		note.Tags[j] = strings.TrimSpace(note.Tags[j])
	}
	return
}
