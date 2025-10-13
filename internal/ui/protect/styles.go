package protect

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	emptyStyle        = lipgloss.NewStyle()
	trashStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	focusedStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("050"))
	blurredStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	unpersistedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("999"))
	warningStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	cursorStyle       = focusedStyle
	activeFilterStyle = lipgloss.NewStyle().Underline(true).Foreground(lipgloss.Color("1"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)
