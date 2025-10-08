package form

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"milesq.dev/btrbk-manage/pkg/components"
)

type Model struct {
	Inputs     []textinput.Model
	focusIndex int
	styles     FormStyles
}

type FormStyles struct {
	BlurredButton string
	FocusedButton string
	BlurStyle     lipgloss.Style
	FocuseStyle   lipgloss.Style
}

func New(inputs []textinput.Model, styles FormStyles) Model {
	return Model{
		Inputs: inputs,
		styles: styles,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd, *components.UpdateMeta) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, nil, &components.UpdateMeta{PassThrough: true}
		case "esc":
			return m, nil, &components.UpdateMeta{Finish: true}
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.Inputs) {
				return m, nil, &components.UpdateMeta{Finish: true}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.Inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.Inputs)
			}

			var cmd tea.Cmd
			for i := range m.Inputs {
				if i == m.focusIndex {
					cmd = m.Inputs[i].Focus()
					m.Inputs[i].PromptStyle = m.styles.FocuseStyle
					m.Inputs[i].TextStyle = m.styles.FocuseStyle
				} else {
					m.Inputs[i].Blur()
					m.Inputs[i].PromptStyle = m.styles.BlurStyle
					m.Inputs[i].TextStyle = m.styles.BlurStyle
				}
			}

			return m, cmd, nil
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd, nil
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &m.styles.BlurredButton
	if m.focusIndex == len(m.Inputs) {
		button = &m.styles.FocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}
