package telescope

import (
	"strings"

	"github.com/Cybernetics354/mayohttp/app/ui"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	search := m.search.View()
	separator := lipgloss.NewStyle().
		Foreground(ui.FocusColor).
		Render(strings.Repeat("-", m.list.Width()))
	list := m.list.View()
	help := m.renderHelp()
	header := lipgloss.NewStyle().Padding(0, 1).Render(m.title)
	description := m.renderDescription()

	telescope := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.FocusColor).
		Render(lipgloss.JoinVertical(lipgloss.Top, search, separator, list, separator, help))

	return lipgloss.JoinVertical(lipgloss.Top, ui.RenderWithHeader(telescope, header), description)
}

func (m Model) renderHelp() string {
	items := []string{
		"↓/<c-j> Down",
		"↑/<c-k> Up",
		"<c-d> Clear",
		"enter Select",
	}

	str := strings.Join(items, " • ")

	return lipgloss.NewStyle().Foreground(ui.BlurColor).Render(str)
}

func (m Model) renderDescription() string {
	item, ok := m.list.SelectedItem().(ui.ListItem)
	if !ok {
		return ""
	}

	if strings.TrimSpace(item.Description()) == "" {
		return ""
	}

	return lipgloss.NewStyle().
		Width(m.width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.FocusColor).
		Render(item.Description())
}
