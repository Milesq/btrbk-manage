package router

import tea "github.com/charmbracelet/bubbletea"

type UpdateHandler[M any] func(tea.Msg) (M, tea.Cmd, *UpdateMeta)

type UpdateRouter[M any] struct {
	model    M
	handlers []routeEntry[M]
}

type routeEntry[M any] struct {
	condition bool
	handler   UpdateHandler[M]
}

func NewRouter[M any](model M) *UpdateRouter[M] {
	return &UpdateRouter[M]{
		model:    model,
		handlers: make([]routeEntry[M], 0),
	}
}

func (r *UpdateRouter[M]) When(condition bool, handler UpdateHandler[M]) *UpdateRouter[M] {
	r.handlers = append(r.handlers, routeEntry[M]{
		condition: condition,
		handler:   handler,
	})
	return r
}

func (r *UpdateRouter[M]) Default(handler UpdateHandler[M]) *UpdateRouter[M] {
	return r.When(true, handler)
}

func (r *UpdateRouter[M]) Update(msg tea.Msg) (M, tea.Cmd) {
	currentModel := r.model
	var combinedCmds []tea.Cmd

	for _, entry := range r.handlers {
		if !entry.condition {
			continue
		}

		newModel, cmd, meta := entry.handler(msg)
		currentModel = newModel

		if cmd != nil {
			combinedCmds = append(combinedCmds, cmd)
		}

		if meta == nil || !meta.PassThrough {
			return currentModel, tea.Batch(combinedCmds...)
		}
	}

	return currentModel, tea.Batch(combinedCmds...)
}
