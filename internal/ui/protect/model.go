package protect

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"milesq.dev/btrbk-manage/internal/app"
	"milesq.dev/btrbk-manage/internal/snaps"
	"milesq.dev/btrbk-manage/pkg/form"
	"milesq.dev/btrbk-manage/pkg/multiselect"
	"milesq.dev/btrbk-manage/pkg/router"
)

type Model struct {
	// Core
	mng     snaps.BackupManager
	dir     string
	backups []snaps.Backup
	cfg     *app.Config

	// General State
	Err         error
	successMsg  string
	cursor      int
	selected    snaps.Backup
	subvolNames []string
	width       int

	// Modes flags
	listProtectedOnly              bool
	trashMode                      bool
	isEdit                         bool
	isConfirmingDelete             bool
	isChoosingSubvolumesForRestore bool

	// SubComponents
	form            form.Model
	restoreSelector multiselect.Model
}

func InitialModel(cfg *app.Config) (Model, error) {
	backupManager := snaps.GetBackupManager(cfg.Paths)
	inputs := getProtectionNoteInputs()

	m := Model{
		dir: cfg.Paths.Snaps,
		mng: backupManager,
		cfg: cfg,
		form: form.New(inputs, form.NewFormProps().WithStyles(form.FormStyles{
			BlurredButton: blurredButton,
			FocusedButton: focusedButton,
			BlurStyle:     lipgloss.NewStyle(),
			FocuseStyle:   focusedStyle,
		})),
	}
	m.recollect()
	return m, nil
}

func (m Model) Init() tea.Cmd {
	return m.form.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
	}

	updatedModel, cmd := router.NewRouter(m).
		When(m.isConfirmingDelete, m.handleDeleteConfirmation).
		When(m.isEdit, m.handleForm).
		When(m.isChoosingSubvolumesForRestore, m.handleRestoreSelector).
		Default(m.handleList).
		Update(msg)
	return updatedModel, cmd
}

func (m *Model) recollect() {
	m.mng.ClearCache()
	backups, err := m.mng.Collect()
	m.Err = err

	filtered := []snaps.Backup{}

	m.subvolNames = backups.SubvolNames

	for _, backup := range backups.Backups {
		if m.trashMode && !backup.IsTrashed {
			continue
		}

		if m.listProtectedOnly && !backup.IsProtected {
			continue
		}

		filtered = append(filtered, backup)
	}

	m.backups = filtered
}
