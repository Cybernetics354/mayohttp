package urlcompose

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/Cybernetics354/mayohttp/app/ui"
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
	helperStyle := lipgloss.NewStyle().
		Background(ui.FocusColor).
		Align(lipgloss.Center)

	var helper []string
	if m.protocol != "" || len(m.paths) > 0 {
		if m.protocol != "" {
			comp := helperStyle.Width(len(m.protocol)).MarginRight(2).Render("0")
			helper = append(helper, comp)
		}

		for i, path := range m.paths {
			index := i + 1
			comp := helperStyle.
				Width(len(path)).
				Render(fmt.Sprintf("%d", index))
			helper = append(helper, comp)
		}
	}

	view := container.Width(m.width).
		Render(lipgloss.JoinVertical(lipgloss.Left, m.url, strings.Join(helper, " ")))

	return ui.RenderWithHeader(view, header)
}

func (m *Model) RenderInput() string {
	return container.Width(m.width).Render(m.input.View())
}
