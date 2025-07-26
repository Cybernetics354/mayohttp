package urlcompose

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) SetWidth(width int) {
	m.width = width
	m.input.Width = width - 3
}

func (m Model) RunCommand() (Model, tea.Cmd) {
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
			err = m.cd(arg)
		}
	}

	if err != nil {
		return m, sendMsg(err)
	}

	m.RefreshUrl()
	m.ClearInput()

	return m, sendMsg(Changed{Url: m.url})
}

func (m *Model) ComposeSuggestions() {
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

func (m *Model) RefreshUrl() string {
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

func (m *Model) ClearInput() {
	m.input.SetValue("")
}

func (m *Model) SetUrl(url string) {
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

func (m *Model) rm(arg string) error {
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

func (m *Model) cd(arg string) error {
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

	return nil
}

func (m *Model) set(command string) error {
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
