package component

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Cybernetics354/mayohttp/app/ui"
)

type UrlCompose struct {
	/// cache the result of the url
	url string

	paths    []string
	queries  map[string]string
	protocol string
	input    textinput.Model

	width int
}

type UrlComposeSubmitMsg struct {
	Url string
}

func NewUrlCompose() UrlCompose {
	m := UrlCompose{
		url:      "",
		paths:    []string{},
		queries:  make(map[string]string),
		protocol: "",
		input:    textinput.New(),
		width:    60,
	}

	m.input.Focus()
	m.input.ShowSuggestions = true
	m.input.SetWidth(m.width)

	return m
}

func (m UrlCompose) Init() tea.Cmd {
	return nil
}

func (m UrlCompose) Update(msg tea.Msg) (UrlCompose, tea.Cmd) {
	var cmd tea.Cmd

	m.input, cmd = m.input.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m, cmd = m.RunCommand()
		case "ctrl+d":
			m.ClearInput()
		}
	}

	return m, cmd
}

func (m *UrlCompose) SetWidth(width int) {
	m.width = width
	m.input.SetWidth(width - 3)
}

func (m UrlCompose) RunCommand() (UrlCompose, tea.Cmd) {
	val := m.input.Value()
	command := strings.SplitN(val, " ", 2)
	var err error

	switch len(command) {
	case 1:
		err = m.set(command[0])
	case 2:
		cmdType, arg := command[0], command[1]
		switch cmdType {
		case "rm":
			err = m.rm(arg)
		case "cd":
			m.cd(arg)
		}
	}

	if err != nil {
		return m, sendMsg(err)
	}

	m.RefreshUrl()
	m.ClearInput()

	return m, sendMsg(UrlComposeSubmitMsg{Url: m.url})
}

func (m *UrlCompose) ComposeSuggestions() {
	suggestions := []string{
		"cd",
		"rm",
		"rm *",
	}

	deleteSug := "rm 0"
	changeSug := fmt.Sprintf("0=%s", m.protocol)
	suggestions = append(suggestions, deleteSug, changeSug)

	for i, path := range m.paths {
		index := i + 1

		deleteSug = fmt.Sprintf("rm %d", index)
		changeSug = fmt.Sprintf("%d=%s", index, path)
		suggestions = append(suggestions, deleteSug, changeSug)
	}

	for k, v := range m.queries {
		deleteSug = fmt.Sprintf("rm %s", k)
		changeSug = fmt.Sprintf("%s=%s", k, v)
		suggestions = append(suggestions, deleteSug, changeSug)
	}

	m.input.SetSuggestions(suggestions)
}

func (m *UrlCompose) RefreshUrl() string {
	defer m.ComposeSuggestions()

	protocol := m.protocol
	if len(protocol) > 0 {
		protocol = fmt.Sprintf("%s://", protocol)
	}

	path := strings.Join(m.paths, "/")
	query := ""

	var queries []string
	for k, v := range m.queries {
		queries = append(queries, fmt.Sprintf("%s=%s", k, v))
	}

	if len(queries) > 0 {
		query = fmt.Sprintf("?%s", strings.Join(queries, "&"))
	}

	m.url = fmt.Sprintf("%s%s%s", protocol, path, query)
	return m.url
}

func (m *UrlCompose) ClearInput() {
	m.input.SetValue("")
}

func (m *UrlCompose) SetUrl(url string) {
	defer m.ComposeSuggestions()

	/// reset the properties
	m.paths = []string{}
	m.queries = make(map[string]string)
	m.protocol = ""

	m.url = url

	parsed := strings.SplitN(url, "?", 2)

	if len(parsed) <= 0 {
		return
	}

	path := parsed[0]
	pathsplit := strings.SplitN(path, "://", 2)

	if len(pathsplit) >= 2 {
		m.protocol = pathsplit[0]
		path = pathsplit[1]
	}

	m.paths = strings.Split(path, "/")

	if len(parsed) < 2 {
		return
	}

	query := parsed[1]
	for q := range strings.SplitSeq(query, "&") {
		queryitem := strings.SplitN(q, "=", 2)
		if len(queryitem) < 2 {
			continue
		}

		m.queries[queryitem[0]] = queryitem[1]
	}
}

func (m *UrlCompose) rm(arg string) error {
	// if not a number, then it should be a query param
	// remove the key from the query
	num, err := strconv.Atoi(arg)
	if err != nil {
		// clear the query params if the arg is *
		if arg == "*" {
			m.queries = make(map[string]string)
			return nil
		}

		delete(m.queries, arg)
		return nil
	}

	// remove the protocol if the number is 0
	if num == 0 {
		m.protocol = ""
		return nil
	}

	if num > len(m.paths) {
		return errors.New("index out of range")
	}

	m.paths = slices.Delete(m.paths, num-1, num)
	return nil
}

func (m *UrlCompose) cd(arg string) {
	// if start with /, then reset the path first
	if strings.HasPrefix(arg, "/") {
		m.paths = []string{}
	}

	for path := range strings.SplitSeq(arg, "/") {
		trimmed := strings.TrimSpace(path)

		switch trimmed {
		case "":
			continue
		case "..":
			if len(m.paths) <= 0 {
				continue
			}
			m.paths = m.paths[:len(m.paths)-1]
			continue
		default:
			m.paths = append(m.paths, trimmed)
			continue
		}
	}
}

func (m *UrlCompose) set(command string) error {
	trimmed := strings.TrimSpace(command)
	parts := strings.SplitN(trimmed, "=", 2)

	if len(parts) != 2 {
		return errors.New("invalid command")
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if num, err := strconv.Atoi(key); err == nil {
		if num > len(m.paths) || num < 0 {
			return errors.New("index out of range")
		}

		if num == 0 {
			m.protocol = value
			return nil
		}

		m.paths[num-1] = value
		return nil
	}

	m.queries[key] = value
	return nil
}

func sendMsg(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func (m UrlCompose) Container() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(ui.FocusColor)
}

func (m UrlCompose) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.RenderUrl(),
		m.RenderInput(),
	)
}

func (m *UrlCompose) RenderUrl() string {
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

	view := m.Container().Width(m.width).
		Render(lipgloss.JoinVertical(lipgloss.Left, m.url, strings.Join(helper, " ")))

	return ui.RenderWithHeader(view, header)
}

func (m *UrlCompose) RenderInput() string {
	return m.Container().Width(m.width).Render(m.input.View())
}
