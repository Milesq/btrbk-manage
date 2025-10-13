package protect

import "github.com/charmbracelet/bubbles/textinput"

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
