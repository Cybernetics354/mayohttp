package component

import (
	"fmt"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/Cybernetics354/mayohttp/app/ui"
)

type UrlInput struct {
	Model  textinput.Model
	method string
	width  int
}

func NewUrlInput(method string) UrlInput {
	return UrlInput{
		Model: ui.UrlInput(method, ""),
	}
}

func (m *UrlInput) SetWidth(width int) *UrlInput {
	m.width = width
	m.recalculateSize()
	return m
}

func (m *UrlInput) SetMethod(method string) *UrlInput {
	m.method = method
	m.Model.Prompt = fmt.Sprintf("%s | ", method)
	m.recalculateSize()
	return m
}

func (m *UrlInput) Value() string {
	return m.Model.Value()
}

func (m *UrlInput) SetValue(value string) *UrlInput {
	m.Model.SetValue(value)
	return m
}

func (m *UrlInput) Clear() tea.Cmd {
	m.Model.SetValue("")
	return nil
}

func (m *UrlInput) Focus() tea.Cmd {
	return m.Model.Focus()
}

func (m *UrlInput) Blur() *UrlInput {
	m.Model.Blur()
	return m
}

func (m UrlInput) Update(msg tea.Msg) (UrlInput, tea.Cmd) {
	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)

	return m, cmd
}

func (m *UrlInput) View() string {
	focused := m.Model.Focused()
	base := m.Model.View()

	if focused {
		return ui.FocusInputContainer.Render(base)
	}

	return ui.BlurInputContainer.Render(base)
}

func (m *UrlInput) recalculateSize() {
	calculatedWidth := m.width - 5 - len(m.method)
	m.Model.SetWidth(calculatedWidth)
}
