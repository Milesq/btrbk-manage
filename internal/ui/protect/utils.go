package protect

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"milesq.dev/btrbk-manage/internal/snaps"
)

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

const building_params_n int = 3

func (m *Model) populateFormWithNote(note snaps.ProtectionNote) {
	if len(m.form.Inputs) < 3 {
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
