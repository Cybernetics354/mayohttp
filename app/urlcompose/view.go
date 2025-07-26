package urlcompose

import (
	"github.com/Cybernetics354/mayohttp/app/ui"
	"github.com/charmbracelet/lipgloss"
)

var container = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(ui.FocusColor)

func (m Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.RenderUrl(),
		m.RenderInput(),
	)
}

func (m *Model) RenderUrl() string {
	header := lipgloss.NewStyle().Padding(0, 1).Render("URL Result")
	return ui.RenderWithHeader(container.Width(m.width).Render(m.url), header)
}

func (m *Model) RenderInput() string {
	return container.Width(m.width).Render(m.input.View())
}
