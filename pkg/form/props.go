package form

import "github.com/charmbracelet/lipgloss"

type FormStyles struct {
	BlurredButton string
	FocusedButton string
	BlurStyle     lipgloss.Style
	FocuseStyle   lipgloss.Style
}

type FormProps struct {
	styles      FormStyles
	clearOnExit bool
}

func NewFormProps() *FormProps {
	return &FormProps{
		clearOnExit: true,
	}
}

func (p *FormProps) WithStyles(styles FormStyles) *FormProps {
	p.styles = styles
	return p
}

func (p *FormProps) WithClearOnExit(clear bool) *FormProps {
	p.clearOnExit = clear
	return p
}

func (p *FormProps) Styles() FormStyles {
	return p.styles
}

func (p *FormProps) ClearOnExit() bool {
	return p.clearOnExit
}
