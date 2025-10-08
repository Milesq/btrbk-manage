package form

import tea "github.com/charmbracelet/bubbletea"

type ExitReason int

const (
	UserSaved ExitReason = iota
	UserCanceled
)

type ExitMsg struct {
	Reason ExitReason
	Values []string
}

func finished(reason ExitReason, values []string) tea.Cmd {
	return func() tea.Msg {
		return ExitMsg{reason, values}
	}
}
