package multiselect

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"milesq.dev/btrbk-manage/pkg/router"
)

type ExitReason int

const (
	UserSaved ExitReason = iota
	UserCanceled
)

type ExitMsg struct {
	Reason   ExitReason
	Selected []string
}

type Styles struct {
	Title        lipgloss.Style
	Item         lipgloss.Style
	SelectedItem lipgloss.Style
	Cursor       lipgloss.Style
	CheckedBox   string
	UncheckedBox string
}

func DefaultStyles() Styles {
	return Styles{
		Title:        lipgloss.NewStyle().Bold(true).MarginBottom(1),
		Item:         lipgloss.NewStyle(),
		SelectedItem: lipgloss.NewStyle().Foreground(lipgloss.Color("212")),
		Cursor:       lipgloss.NewStyle().Foreground(lipgloss.Color("212")),
		CheckedBox:   "[x]",
		UncheckedBox: "[ ]",
	}
}

type Model struct {
	title       string
	options     []string
	selected    map[string]bool
	cursor      int
	styles      Styles
	preselected []string
}

func New(title string, options []string, preselected []string, styles *Styles) Model {
	s := DefaultStyles()
	if styles != nil {
		s = *styles
	}

	selected := make(map[string]bool)
	for _, opt := range preselected {
		selected[opt] = true
	}

	return Model{
		title:       title,
		options:     options,
		selected:    selected,
		styles:      s,
		preselected: preselected,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd, *router.UpdateMeta) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, finished(UserCanceled, nil), router.PassThrough()
		case "esc":
			return m, finished(UserCanceled, nil), nil
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.options) - 1
			}
		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.options) {
				m.cursor = 0
			}
		case " ":
			if len(m.options) > 0 {
				opt := m.options[m.cursor]
				m.selected[opt] = !m.selected[opt]
			}
		case "enter":
			return m, finished(UserSaved, m.GetSelected()), nil
		case "a":
			for _, opt := range m.options {
				m.selected[opt] = true
			}
		case "n":
			for _, opt := range m.options {
				m.selected[opt] = false
			}
		}
	}
	return m, nil, nil
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(m.styles.Title.Render(m.title))
	b.WriteString("\n")

	for i, opt := range m.options {
		cursor := "  "
		if i == m.cursor {
			cursor = m.styles.Cursor.Render("> ")
		}

		checkbox := m.styles.UncheckedBox
		if m.selected[opt] {
			checkbox = m.styles.CheckedBox
		}

		itemText := checkbox + " " + opt
		if m.selected[opt] {
			b.WriteString(cursor + m.styles.SelectedItem.Render(itemText))
		} else {
			b.WriteString(cursor + m.styles.Item.Render(itemText))
		}
		b.WriteString("\n")
	}

	b.WriteString("\nPress space to toggle, a=all, n=none, enter to submit, esc to cancel\n")

	return b.String()
}

func (m Model) GetSelected() []string {
	var result []string
	for _, opt := range m.options {
		if m.selected[opt] {
			result = append(result, opt)
		}
	}
	return result
}

func finished(reason ExitReason, selected []string) tea.Cmd {
	return func() tea.Msg {
		return ExitMsg{reason, selected}
	}
}
