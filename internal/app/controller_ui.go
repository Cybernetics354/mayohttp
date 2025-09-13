package app

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	c "github.com/Cybernetics354/mayohttp/internal/app/common"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/event"
	"github.com/Cybernetics354/mayohttp/internal/app/ui/telescope"
	"github.com/Cybernetics354/mayohttp/internal/app/val"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *State) ClearFocusedInput() (tea.Model, tea.Cmd) {
	field := m.GetFocusedField()
	if field == nil {
		return m, nil
	}

	switch f := field.(type) {
	case *textinput.Model:
		f.SetValue("")
	}

	return m, nil
}

func (m *State) OpenTelescope(msg event.OpenTelescope) (tea.Model, tea.Cmd) {
	var items []list.Item
	var title string

	switch msg.TeleType {
	case val.TELESCOPE_METHOD_PALLETE:
		items = val.MethodPalletesTelescope
		title = "Select Method"
	case val.TELESCOPE_QUICK_ACCESS:
		items = val.QuickAccess
		title = "Quick Access"
	}

	m.telescope.SetList(items)
	m.telescope.SetTitle(title)
	m.telescope.SetTeleType(msg.TeleType)

	return m, c.SendMsg(event.AddStack{State: val.STATE_TELESCOPE})
}

func (m *State) SelectTelescopeItem(msg telescope.SubmitMsg) (tea.Model, tea.Cmd) {
	switch msg.TeleType {
	case val.TELESCOPE_QUICK_ACCESS:
		v, ok := msg.Value.Value().([]tea.Msg)
		if !ok {
			return m, c.SendMsg(event.Err(errors.New("no menu selected")))
		}

		var cmds []tea.Cmd
		for _, msg := range v {
			cmds = append(cmds, c.SendMsg(msg))
		}

		return m, tea.Sequence(c.SendMsg(event.PopStack{}), tea.Batch(cmds...))
	case val.TELESCOPE_METHOD_PALLETE:
		v, ok := msg.Value.Value().(string)
		if !ok {
			return m, c.SendMsg(event.Err(errors.New("no method selected")))
		}

		m.telescope.Clear()
		m.method = v
		m.url.Prompt = v + " | "
		m.url.Width = m.sw - 5 - len(m.url.Prompt)
	}

	return m, c.SendMsg(event.PopStack{})
}

func (m *State) CopyToClipboard() (tea.Model, tea.Cmd) {
	var v string
	switch m.state {
	case val.STATE_FOCUS_URL:
		v = m.url.Value()
	case val.STATE_FOCUS_PIPE:
		v = m.pipe.Value()
	case val.STATE_FOCUS_PIPEDRESP:
		v = m.pipedresp.Value()
	}

	err := clipboard.WriteAll(v)
	if err != nil {
		return m, tea.Batch(
			c.SendMsg(event.Err(err)),
			c.SendMsg(event.SetActivity("Error copying to clipboard")),
		)
	}
	return m, c.SendMsg(event.SetActivity("Copied to clipboard"))
}

func (m *State) HandleListFilter(msg list.FilterMatchesMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case val.STATE_COMMAND_PALLETE:
		m.commands, cmd = m.commands.Update(msg)
	case val.STATE_METHOD_PALLETE:
		m.methodSelect, cmd = m.methodSelect.Update(msg)
	case val.STATE_SELECT_ENV:
		m.envList, cmd = m.envList.Update(msg)
	case val.STATE_SELECT_SESSION, val.STATE_SAVE_SESSION:
		m.sessionList, cmd = m.sessionList.Update(msg)
	}

	return m, cmd
}

func (m *State) GetField(state string) any {
	switch state {
	case val.STATE_FOCUS_URL:
		return &m.url
	case val.STATE_FOCUS_PIPE:
		return &m.pipe
	case val.STATE_FOCUS_PIPEDRESP:
		return &m.pipedresp
	case val.STATE_FOCUS_RESPONSE:
		return &m.response
	case val.STATE_FOCUS_BODY:
		return &m.body
	case val.STATE_FOCUS_HEADER:
		return &m.header
	case val.STATE_COMMAND_PALLETE:
		return &m.commands
	case val.STATE_FOCUS_RESPONSE_FILTER:
		return &m.resFilter
	case val.STATE_SAVE_SESSION_INPUT, val.STATE_SESSION_RENAME_INPUT:
		return &m.saveInput
	}

	return nil
}

func (m *State) GetFocusedField() any {
	return m.GetField(m.state)
}

func (m *State) OpenEditor(msg event.OpenEditor) (tea.Model, tea.Cmd) {
	var str string
	switch f := m.GetField(msg.State).(type) {
	case *textarea.Model:
		str = f.Value()
	case *textinput.Model:
		str = f.Value()
	default:
		return m, nil
	}

	dir := filepath.Dir(val.TempFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return m, c.SendMsg(event.Err(err))
		}
	}

	f, err := os.Create(val.TempFilePath)
	if err != nil {
		return m, c.SendMsg(event.Err(err))
	}

	defer f.Close()
	if _, err := f.WriteString(str); err != nil {
		return m, c.SendMsg(event.Err(err))
	}

	editor := c.GetDefaultEditor()

	cmd := exec.Command(editor, val.TempFilePath)
	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		if err != nil {
			return event.Err(err)
		}

		f, err := os.Open(val.TempFilePath)
		if err != nil {
			return event.Err(err)
		}

		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return event.Err(err)
		}

		str := strings.TrimSpace(string(b))

		return event.SetFieldValue{State: msg.State, Value: str}
	})
}

func (m *State) SetFieldValue(msg event.SetFieldValue) (tea.Model, tea.Cmd) {
	switch f := m.GetField(msg.State).(type) {
	case *textarea.Model:
		f.SetValue(msg.Value)
	case *textinput.Model:
		f.SetValue(msg.Value)
	}

	return m, c.SendMsg(event.SetActivity("Set Field Value"))
}

func (m *State) ShowSpinner() (tea.Model, tea.Cmd) {
	m.showSpinner = true
	return m, nil
}

func (m *State) HideSpinner() (tea.Model, tea.Cmd) {
	m.showSpinner = false
	return m, nil
}

func (m *State) SelectCommandPallete() (tea.Model, tea.Cmd) {
	if m.commands.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.commands, cmd = m.commands.Update(tea.KeyMsg{
			Type: tea.KeyEnter,
		})
		return m, cmd
	}

	i, ok := m.commands.SelectedItem().(val.CommandPallete)
	if !ok {
		return m, c.SendMsg(event.Err(errors.New("no command selected")))
	}

	return m, c.SendMsg(event.RunCommand{CommandId: i.CommandId})
}

func (m *State) SelectMethodPallete() (tea.Model, tea.Cmd) {
	if m.methodSelect.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.methodSelect, cmd = m.methodSelect.Update(tea.KeyMsg{
			Type: tea.KeyEnter,
		})
		return m, cmd
	}

	i, ok := m.methodSelect.SelectedItem().(val.MethodPallete)
	if !ok {
		return m, c.SendMsg(event.Err(errors.New("no method selected")))
	}

	m.method = i.Method
	m.url.Prompt = i.Method + " | "
	m.url.Width = m.sw - 5 - len(m.url.Prompt)

	return m, c.SendMsg(event.PopStack{})
}

func (m *State) RunCommand(command event.RunCommand) (tea.Model, tea.Cmd) {
	switch command.CommandId {
	case val.COMMAND_OPEN_ENV:
		return m, c.SendMsg(event.OpenEnv{})
	case val.COMMAND_SELECT_METHOD:
		return m, c.SendMsg(event.OpenTelescope{TeleType: val.TELESCOPE_METHOD_PALLETE})
	case val.COMMAND_OPEN_BODY:
		return m, c.SendMsg(event.OpenRequestBody{})
	case val.COMMAND_OPEN_HEADER:
		return m, c.SendMsg(event.OpenRequestHeader{})
	case val.COMMAND_SAVE_SESSION:
		return m, tea.Sequence(
			c.SendMsg(event.AddStack{State: val.STATE_SAVE_SESSION}),
			c.SendMsg(event.LoadSessionList{}),
		)
	case val.COMMAND_OPEN_SESSION_LIST:
		return m, tea.Sequence(
			c.SendMsg(event.AddStack{State: val.STATE_SELECT_SESSION}),
			c.SendMsg(event.LoadSessionList{}),
		)
	case val.COMMAND_CHANGE_ENV:
		return m, tea.Batch(
			c.SendMsg(event.RefreshSelectEnv{}),
			c.SendMsg(event.AddStack{State: val.STATE_SELECT_ENV}),
		)

	default:
		return m, nil
	}
}

func (m *State) OpenRequestBody() (tea.Model, tea.Cmd) {
	return m, c.SendMsg(event.OpenEditor{State: val.STATE_FOCUS_BODY})
}

func (m *State) OpenRequestHeader() (tea.Model, tea.Cmd) {
	return m, c.SendMsg(event.OpenEditor{State: val.STATE_FOCUS_HEADER})
}
