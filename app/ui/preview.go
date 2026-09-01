package ui

import "charm.land/lipgloss/v2"

type Preview struct {
	Header    string
	Width     int
	MaxHeight int
	Body      string
}

func (p *Preview) Render() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		PreviewHeaderStyle.Render(p.Header),
		PreviewBodyStyle.Width(p.Width).
			Render(lipgloss.NewStyle().MaxHeight(p.MaxHeight).Render(p.Body)),
	)
}
