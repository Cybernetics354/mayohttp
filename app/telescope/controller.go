package telescope

import (
	"github.com/Cybernetics354/mayohttp/app/ui"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func sendMsg(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func (t *Model) Blur() {
	t.search.Blur()
}

func (t *Model) Focus() {
	t.search.Focus()
}

func (t *Model) Clear() {
	t.searchRegistry[t.teleType] = ""
	t.search.SetValue("")
}

func (t *Model) Sync() {
	if t.list.FilterValue() != t.search.Value() {
		t.list.SetFilterText(t.search.Value())
	}

	t.searchRegistry[t.teleType] = t.search.Value()
}

func (t *Model) RefreshRender() {
	t.search.Width = t.width - 3
	t.list.SetSize(t.width, t.height)
}

func (t *Model) SetList(list []list.Item) {
	t.list.SetItems(list)
	t.list.SetFilterText(t.search.Value())
}

func (t *Model) SetTitle(title string) {
	t.title = title
}

func (t *Model) SetSize(width, height int) {
	t.width = width
	t.height = height

	t.RefreshRender()
}

func (t *Model) SetTeleType(teleType any) {
	if t.teleType != teleType {
		t.search.SetValue(t.searchRegistry[teleType])
	}

	t.teleType = teleType
	t.Sync()
}

func (t *Model) GetSelectedMsg() (cmd tea.Cmd) {
	item, ok := t.list.SelectedItem().(ui.ListItem)
	if !ok {
		return sendMsg(ErrorMsg("no item selected"))
	}

	return sendMsg(SubmitMsg{
		Title:    t.title,
		TeleType: t.teleType,
		Value:    item,
	})
}
