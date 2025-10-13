package protect

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"milesq.dev/btrbk-manage/internal/snaps"
	"milesq.dev/btrbk-manage/pkg/form"
	"milesq.dev/btrbk-manage/pkg/router"
)

type Model struct {
	// Core
	mng     snaps.BackupManager
	dir     string
	backups []snaps.Backup

	// General State
	Err      error
	cursor   int
	selected snaps.Backup

	// Modes flags
	listProtectedOnly  bool
	trashMode          bool
	isEdit             bool
	isConfirmingDelete bool

	// SubComponents
	form form.Model
}

func InitialModel(dir string) Model {
	backupManager := snaps.GetManagerForDirectory(dir)
	info, err := backupManager.Collect()

	inputs := getProtectionNoteInputs()

	return Model{
		backups: info.Backups,
		Err:     err,
		dir:     dir,
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
	updatedModel, cmd := router.NewRouter(m).
		When(m.isConfirmingDelete, m.handleDeleteConfirmation).
		When(m.isEdit, m.handleForm).
		Default(m.handleList).
		Update(msg)
	return updatedModel, cmd
}

func (m *Model) recollect() {
	m.mng.ClearCache()
	backups, err := m.mng.Collect()
	m.Err = err
	m.backups = backups.Backups
}
