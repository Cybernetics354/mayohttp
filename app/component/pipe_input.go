package component

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/Cybernetics354/mayohttp/app/ui"
)

type PipeInput struct {
	Model textinput.Model
	width int
}

func NewPipeInput() PipeInput {
	return PipeInput{
		Model: ui.PipeInput(),
	}
}

func (m *PipeInput) SetWidth(width int) *PipeInput {
	m.width = width
	m.recalculateSize()
	return m
}

func (m *PipeInput) Value() string {
	return m.Model.Value()
}

func (m *PipeInput) SetValue(value string) *PipeInput {
	m.Model.SetValue(value)
	return m
}

func (m *PipeInput) Clear() tea.Cmd {
	m.Model.SetValue("")
	return nil
}

func (m *PipeInput) Focus() tea.Cmd {
	return m.Model.Focus()
}

func (m *PipeInput) Blur() *PipeInput {
	m.Model.Blur()
	return m
}

func (m PipeInput) Update(msg tea.Msg) (PipeInput, tea.Cmd) {
	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)

	return m, cmd
}

func (m *PipeInput) View() string {
	focused := m.Model.Focused()
	base := m.Model.View()

	if focused {
		return ui.FocusInputContainer.Render(base)
	}

	return ui.BlurInputContainer.Render(base)
}

func (m *PipeInput) recalculateSize() {
	calculatedWidth := m.width - 11
	m.Model.SetWidth(calculatedWidth)
}
