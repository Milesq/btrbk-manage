package protect

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"milesq.dev/btrbk-manage/internal/utils"
)

func (m Model) View() string {
	if m.Err != nil {
		return fmt.Sprintf("Error: %v\n\nDir: %s\nPress q to quit.\n", m.Err, m.Dir)
	}
	var b strings.Builder
	title := fmt.Sprintf("Btrbk backups in %s  —  %d backups\n", m.Dir, len(m.Backups))
	b.WriteString(title)
	b.WriteString(strings.Repeat("─", utils.MinMax(10, len(title), 80)))
	b.WriteString("\n\n")

	if m.IsConfirmingDelete {
		m.viewDeleteConfirmation(&b)
	} else if m.IsEdit {
		b.WriteString(m.form.View())
	} else {
		m.viewList(&b)
	}

	return b.String()
}

func (m Model) viewList(b *strings.Builder) {
	if len(m.Backups) == 0 {
		b.WriteString("No snapshot groups found.\n")
		return
	}

	for i, g := range m.Backups {
		if i == m.Cursor {
			b.WriteString(focusedStyle.Render("> "))
		} else {
			b.WriteString("  ")
		}

		if g.IsProtected {
			b.WriteString(focusedStyle.Render("★ "))
		} else if g.IsTrashed {
			b.WriteString(trashStyle.Render("🗑 "))
		} else {
			b.WriteString("  ")
		}

		timestampStyle := lipgloss.NewStyle()
		if g.IsProtected {
			timestampStyle = unpersistedStyle
		}

		b.WriteString(timestampStyle.Render(utils.PrettifyDate(g.Timestamp)))

		if g.ProtectionNote.Note != "" || g.ProtectionNote.Reason != "" || len(g.ProtectionNote.Tags) > 0 {
			b.WriteString(" ")
			b.WriteString(blurredStyle.Render("—"))
			if g.ProtectionNote.Note != "" {
				note := g.ProtectionNote.Note
				if len(note) > 30 {
					note = note[:30] + "..."
				}
				b.WriteString(" ")
				b.WriteString(blurredStyle.Render(note))
			}
			if g.ProtectionNote.Reason != "" {
				b.WriteString(" ")
				b.WriteString(blurredStyle.Render("[" + g.ProtectionNote.Reason + "]"))
			}
			if len(g.ProtectionNote.Tags) > 0 {
				b.WriteString(" ")
				b.WriteString(blurredStyle.Render(strings.Join(g.ProtectionNote.Tags, ", ")))
			}
		}

		b.WriteRune('\n')
	}

	writeHelpMessage(b)
}

func getProtectionNoteInputs() []textinput.Model {
	inputs := make([]textinput.Model, 3)

	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 100
		t.Width = t.CharLimit // FIX: required by bug in bubbles lib, without it, the placeholder is not displayed

		switch i {
		case 0:
			t.Prompt = "Note   " + t.Prompt

			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Prompt = "Reason " + t.Prompt
		case 2:
			t.Prompt = "Tags   " + t.Prompt
			t.Placeholder = "(comma separated)"
		}

		inputs[i] = t
	}

	return inputs
}
