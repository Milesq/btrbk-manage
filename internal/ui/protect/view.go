package protect

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"milesq.dev/btrbk-manage/internal/utils"
)

var (
	focusedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("050"))
	blurredStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	unpersistedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("999"))
	cursorStyle      = focusedStyle

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
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

	if m.SelectedForEdit != nil {
		b.WriteString(m.form.View())
	} else {
		m.ViewList(&b)
	}

	return b.String()
}

func (m Model) ViewList(b *strings.Builder) {
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
			b.WriteString(focusedStyle.Render("★"))
		} else {
			b.WriteString(" ")
		}

		timestampStyle := lipgloss.NewStyle()
		if g.IsProtected {
			timestampStyle = unpersistedStyle
		}

		b.WriteString(timestampStyle.Render(utils.PrettifyDate(g.Timestamp)))
		b.WriteRune('\n')
	}

	dot := focusedStyle.Render(" • ")
	fmt.Fprint(
		b,
		"\n↑/↓ ",
		blurredStyle.Render("to move"),
		dot,
		"Space ",
		blurredStyle.Render("to un/protect backup"),
		dot,
		"Enter ",
		blurredStyle.Render("to edit note"),
		dot,
		"q ",
		blurredStyle.Render("to quit\n"),
	)
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
