package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListItem struct {
	title, desc string
	value       any
}

func NewListItem(title, desc string, value any) ListItem {
	return ListItem{
		title: title,
		desc:  desc,
		value: value,
	}
}

func (l ListItem) Title() string {
	return l.title
}

func (l ListItem) Description() string {
	return l.desc
}

func (l ListItem) FilterValue() string {
	return l.title + " " + l.desc
}

func (l ListItem) Value() any {
	return l.value
}

type ListDelegate struct {
	tile      lipgloss.Style
	focusTile lipgloss.Style
	height    int
	spacing   int
	onUpdate  func(msg tea.Msg, m *list.Model) tea.Cmd
}

func NewListDelegate() ListDelegate {
	return ListDelegate{
		tile:      ListItemStyle,
		focusTile: FocusListItemStyle,
		height:    0,
		spacing:   0,
		onUpdate:  nil,
	}
}

func (d ListDelegate) SetTile(tile lipgloss.Style) ListDelegate {
	d.tile = tile
	return d
}

func (d ListDelegate) SetFocusTile(focusTile lipgloss.Style) ListDelegate {
	d.focusTile = focusTile
	return d
}

func (d ListDelegate) SetHeight(height int) ListDelegate {
	d.height = height
	return d
}

func (d ListDelegate) SetSpacing(spacing int) ListDelegate {
	d.spacing = spacing
	return d
}

func (d ListDelegate) SetOnUpdate(onUpdate func(msg tea.Msg, m *list.Model) tea.Cmd) ListDelegate {
	d.onUpdate = onUpdate
	return d
}

func (d ListDelegate) Height() int {
	return 1
}

func (d ListDelegate) Spacing() int {
	return 0
}

func (d ListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	if d.onUpdate == nil {
		return nil
	}

	return d.onUpdate(msg, m)
}

func (d ListDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(ListItem)
	if !ok {
		return
	}

	style := d.tile.Width(m.Width())
	selected := m.Index() == index

	if selected {
		style = d.focusTile.Width(m.Width())
	}

	numeral := fmt.Sprintf("%d. ", index+1)

	fmt.Fprint(
		w,
		style.Render(lipgloss.JoinHorizontal(lipgloss.Left, numeral, item.Title())),
	)
}

func CreateList() list.Model {
	i := list.New([]list.Item{}, NewListDelegate(), 0, 0)
	i.SetShowStatusBar(false)
	i.SetShowTitle(false)
	i.SetShowFilter(false)
	i.SetShowHelp(false)

	return i
}
