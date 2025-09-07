package app

import (
	"fmt"
	"math"
	"slices"

	c "github.com/Cybernetics354/mayohttp/internal/app/common"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/intf"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/store"
	"github.com/Cybernetics354/mayohttp/internal/app/ui"
	"github.com/Cybernetics354/mayohttp/internal/app/val"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

func (m State) View() string {
	return m.Render()
}

func (m *State) RefreshView() {
	w, h := m.sw, m.sh

	lh := h - 1

	m.help.Width = w
	m.url.Width = w - 5 - len(m.url.Prompt)
	m.pipe.Width = w - 11
	m.response.SetWidth(w)
	m.response.SetHeight(h - 10)
	m.pipedresp.SetWidth(w)
	m.pipedresp.SetHeight(h - 9)
	m.commands.SetSize(ui.ListWidth, lh)
	m.envList.SetSize(ui.ListWidth, lh)
	m.sessionList.SetSize(ui.ListWidth, lh)
	m.methodSelect.SetSize(w, h)

	m.telescope.SetSize(c.Clamp(w/2, 60, 90), c.Clamp(h/4, 10, 30))
	m.urlcompose.SetWidth(c.Clamp(w/2, 60, 90))
}

func (m *State) Render() string {
	base := m.RenderBase()
	overlay := m.GetOverlayLayers()
	view := ui.NewCompositeView(base)

	for _, layer := range overlay {
		view.AddLayer(layer)
	}

	return val.AppStyle.Render(view.Render())
}

func (m *State) RenderBase() string {
	state := m.state

	if slices.Contains(val.Overlays, state) {
		length := len(m.stateStack)
		for i := range m.stateStack {
			index := (length - 1) - i
			cState := m.stateStack[index]

			if slices.Contains(val.Overlays, cState) {
				continue
			}

			state = cState
			break
		}
	}

	switch state {
	case val.STATE_COMMAND_PALLETE:
		return m.RenderWithListHelp(val.ListMapping, m.RenderCommandPallete())
	case val.STATE_METHOD_PALLETE:
		return lipgloss.JoinVertical(lipgloss.Top, m.methodSelect.View())
	case val.STATE_SELECT_ENV:
		return m.RenderWithListHelp(val.ListMapping, m.RenderEnvList())
	case val.STATE_SELECT_SESSION, val.STATE_SAVE_SESSION:
		var mapping help.KeyMap
		mapping = val.ListMapping

		switch m.state {
		case val.STATE_SELECT_SESSION:
			mapping = val.SessionListMapping
		case val.STATE_SAVE_SESSION:
			mapping = val.SaveListMapping
		}

		return m.RenderWithListHelp(mapping, m.RenderSessionList())
	default:
		return lipgloss.JoinVertical(
			lipgloss.Top,
			m.RenderURL(),
			m.RenderPipe(),
			lipgloss.NewStyle().PaddingLeft(1).Render(
				m.resFilter.Render(),
			),
			m.RenderPipedResponse(),
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.RenderHelp(),
				"  | ",
				m.RenderSpinner(),
				m.activity,
				" | ",
				fmt.Sprintf("ENV(%s)", val.EnvFilePath),
			),
		)
	}
}

func (m *State) GetOverlayLayers() []*ui.CompositeViewLayer {
	var layers []*ui.CompositeViewLayer

	for _, state := range m.stateStack {
		if !slices.Contains(val.Overlays, state) {
			continue
		}

		layer := ui.NewCompositeViewLayer()
		switch state {
		case val.STATE_KEYBINDING_MODAL:
			layer.SetView(m.RenderKeybindings())
		case val.STATE_SAVE_SESSION_INPUT, val.STATE_SESSION_RENAME_INPUT:
			layer.SetView(m.RenderSessionInput())
		case val.STATE_URL_COMPOSE:
			layer.SetView(m.urlcompose.View())
		case val.STATE_TELESCOPE:
			layer.SetView(m.telescope.View())
			layer.SetPositionY(ui.CompositeLayerTop)
			layer.SetOffset(0, m.sh/3)
		}

		layers = append(layers, layer)
	}

	return layers
}

func (m *State) RenderSessionInput() string {
	title := lipgloss.NewStyle().Foreground(ui.FocusColor).Padding(0, 1).Render(m.saveInput.Prompt)
	m.saveInput.Prompt = ""

	base := lipgloss.NewStyle().
		Width(int(math.Max(60, float64(m.sw/2)))).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ui.FocusColor).
		Render(m.saveInput.View())

	return ui.RenderWithHeader(base, title)
}

func (m *State) RenderKeybindings() string {
	title := lipgloss.NewStyle().Foreground(ui.FocusColor).Padding(0, 1).Render("Keybinding")
	base := lipgloss.NewStyle().
		MaxWidth(int(math.Max(60, float64(m.sw/2)))).
		PaddingRight(2).
		MaxHeight(m.sh - (m.sh/12)*2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ui.FocusColor)

	var tiles []string
	keyStyle := lipgloss.NewStyle().
		Width(12).
		Foreground(ui.FocusColor).
		Align(lipgloss.Right).
		PaddingRight(1)

	tiles = append(
		tiles,
		lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render(""), "-- Local --"),
	)
	for _, key := range val.HomeMapping.KeybindingHelp() {
		key, desc := key.Help().Key, key.Help().Desc

		tile := lipgloss.JoinHorizontal(
			lipgloss.Left,
			keyStyle.Render(key),
			desc,
		)

		tiles = append(tiles, tile)
	}

	content := lipgloss.JoinVertical(lipgloss.Top, tiles...)

	return ui.RenderWithHeader(base.Render(content), title)
}

func (m *State) PreviewSize() (int, int) {
	return m.sw - ui.ListPreviewWidthMargin, m.sh - 7
}

func (m *State) RenderSessionList() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		ui.ListContainer.Render(m.sessionList.View()),
		" ",
		m.RenderSessionListPreview(),
	)
}

func (m *State) RenderSessionListPreview() string {
	pw, ph := m.PreviewSize()
	item, ok := m.sessionList.SelectedItem().(*store.Session)
	if !ok {
		return ""
	}

	str := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, item.Method, " : ", item.Url),
		lipgloss.JoinHorizontal(lipgloss.Top, "PIPE : ", item.Pipe),
		item.PipedResponse,
	)

	preview := ui.Preview{
		Header:    "Preview",
		Width:     pw,
		MaxHeight: ph,
		Body:      str,
	}

	return preview.Render()
}

func (m *State) RenderEnvList() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		ui.ListContainer.Render(m.envList.View()),
		" ",
		m.RenderEnvListPreview(),
	)
}

func (m *State) RenderEnvListPreview() string {
	pw, ph := m.PreviewSize()
	file, ok := m.envList.SelectedItem().(intf.FileItem)

	prev := ""
	if ok {
		uiPrev := ui.Preview{
			Header:    "Preview",
			Width:     pw,
			MaxHeight: ph,
			Body:      c.Printval(file.Path, true),
		}

		prev = lipgloss.JoinVertical(
			lipgloss.Left,
			uiPrev.Render(),
		)
	}

	return prev
}

func (m *State) RenderCommandPallete() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		ui.ListContainer.Render(m.commands.View()),
		" ",
		m.RenderCommandPalletePreview(),
	)
}

func (m *State) RenderCommandPalletePreview() string {
	pw, ph := m.PreviewSize()
	var str string

	command, ok := m.commands.SelectedItem().(val.CommandPallete)
	if !ok {
		return ""
	}

	switch command.CommandId {
	case val.COMMAND_OPEN_ENV:
		str = c.Printval(val.EnvFilePath, true)
	case val.COMMAND_OPEN_BODY:
		str = c.Printval(m.body.Value(), false)
	case val.COMMAND_OPEN_HEADER:
		str = c.Printval(m.header.Value(), false)
	case val.COMMAND_SELECT_METHOD:
		str = fmt.Sprintf("Current method : %s", m.method)
	case val.COMMAND_CHANGE_ENV:
		str = fmt.Sprintf("Current ENV : %s", val.EnvFilePath)
	}

	if len(str) <= 0 {
		return ""
	}

	preview := ui.Preview{
		Header:    "Preview",
		Width:     pw,
		MaxHeight: ph,
		Body:      str,
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		preview.Render(),
	)
}

func (m *State) RenderHelp() string {
	return m.help.View(m.keys)
}

func (m *State) RenderWithListHelp(mapping help.KeyMap, body string) string {
	return lipgloss.JoinVertical(lipgloss.Left, body, m.help.View(mapping))
}

func (m *State) RenderURL() string {
	c := m.url.View()

	if m.state == val.STATE_FOCUS_URL {
		return ui.FocusInputContainer.Render(c)
	}

	return ui.BlurInputContainer.Render(c)
}

func (m *State) RenderPipe() string {
	c := m.pipe.View()

	if m.state == val.STATE_FOCUS_PIPE {
		return ui.FocusInputContainer.Render(c)
	}

	return ui.BlurInputContainer.Render(c)
}

func (m *State) RenderResponse() string {
	return m.response.View()
}

func (m *State) RenderPipedResponse() string {
	return m.pipedresp.View()
}

func (m *State) RenderSpinner() string {
	if !m.showSpinner {
		return ""
	}

	return fmt.Sprintf("%s ", m.spinner.View())
}
