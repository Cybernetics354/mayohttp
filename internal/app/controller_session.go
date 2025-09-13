package app

import (
	"errors"
	"fmt"
	"os"
	"strings"

	c "github.com/Cybernetics354/mayohttp/internal/app/common"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/event"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/store"
	"github.com/Cybernetics354/mayohttp/internal/app/val"
	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *State) applySession(s *store.Session) *State {
	m.url.SetValue(s.Url)
	m.pipe.SetValue(s.Pipe)
	m.pipedresp.SetValue(s.PipedResponse)
	m.response.SetValue(s.Response)
	m.header.SetValue(s.Header)
	m.body.SetValue(s.Body)
	m.method = s.Method
	m.resFilter = s.ResFilter
	m.url.Prompt = m.method + " | "
	m.url.Width = m.sw - 5 - len(m.url.Prompt)

	return m
}

func createSessionFromState(s *State) *store.Session {
	return &store.Session{
		Url:           s.url.Value(),
		Pipe:          s.pipe.Value(),
		PipedResponse: s.pipedresp.Value(),
		Method:        s.method,
		Response:      s.response.Value(),
		Header:        s.header.Value(),
		Body:          s.body.Value(),
		ResFilter:     s.resFilter,
	}
}

func (m *State) ReplaceCurrentSession(msg event.ReplaceCurrentSession) (tea.Model, tea.Cmd) {
	path := msg.Path
	_, err := store.OpenSessionFromPath(path)
	if err != nil {
		return m, tea.Batch(
			c.SendMsg(event.Err(err)),
			c.SendMsg(event.SetActivity("Can't load session on "+path)),
		)
	}

	return m, c.SendMsg(event.LoadSession{Path: path})
}

func (m *State) SaveSession(msg event.SaveSession) (tea.Model, tea.Cmd) {
	session := createSessionFromState(m)
	err := session.Save(msg.Path)
	if err != nil {
		return m, tea.Batch(
			c.SendMsg(event.Err(err)),
			c.SendMsg(event.SetActivity("Error saving session")),
		)
	}

	return m, c.SendMsg(event.SetActivity("Session saved to " + msg.Path))
}

func (m *State) LoadSession(msg event.LoadSession) (tea.Model, tea.Cmd) {
	session, err := store.OpenSessionFromPath(msg.Path)
	if err != nil {
		return m, tea.Batch(
			c.SendMsg(event.Err(err)),
			c.SendMsg(event.SetActivity("Can't load session")),
		)
	}

	m.applySession(session)

	return m, tea.Batch(
		c.SendMsg(event.SetActivity("Session loaded from "+msg.Path)),
		c.SendMsg(event.RecalculateComponentSizes{}),
	)
}

func (m *State) LoadSessionList() (tea.Model, tea.Cmd) {
	if _, err := os.Stat(val.CollectionFolder); os.IsNotExist(err) {
		return m, c.SendMsg(event.Err(err))
	}

	files, err := os.ReadDir(val.CollectionFolder)
	if err != nil {
		return m, c.SendMsg(event.Err(err))
	}

	var sessionItems []list.Item
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		session, err := store.OpenSessionFromPath(
			fmt.Sprintf("%s/%s", val.CollectionFolder, file.Name()),
		)
		if err != nil {
			continue
		}

		session.Name = file.Name()

		sessionItems = append(sessionItems, session)
	}

	m.sessionList.SetItems(sessionItems)
	if m.sessionList.FilterState() == list.FilterApplied {
		m.sessionList.SetFilterText(m.sessionList.FilterValue())
	}

	switch m.state {
	case val.STATE_SAVE_SESSION:
		m.sessionList.Title = "Save Session To"
	case val.STATE_SELECT_SESSION:
		m.sessionList.Title = "Open Session"
	}

	return m, nil
}

func (m *State) DeleteSessionItem() (tea.Model, tea.Cmd) {
	i, ok := m.sessionList.SelectedItem().(*store.Session)
	if !ok {
		return m, c.SendMsg(event.Err(errors.New("no session selected")))
	}

	err := i.Delete()
	if err != nil {
		return m, c.SendMsg(event.Err(err))
	}

	return m, c.SendMsg(event.LoadSessionList{})
}

func (m *State) SelectSessionItem() (tea.Model, tea.Cmd) {
	if m.sessionList.FilterState() == list.Filtering {
		var cmd tea.Cmd
		m.sessionList, cmd = m.sessionList.Update(tea.KeyMsg{
			Type: tea.KeyEnter,
		})
		return m, cmd
	}

	i, ok := m.sessionList.SelectedItem().(*store.Session)
	if !ok {
		return m, c.SendMsg(event.Err(errors.New("no session selected")))
	}

	switch m.state {
	case val.STATE_SAVE_SESSION:
		return m, tea.Batch(
			c.SendMsg(event.SaveSession{Path: i.Path()}),
			c.SendMsg(event.PopStackRoot{}),
		)
	case val.STATE_SELECT_SESSION:
		return m, tea.Batch(
			c.SendMsg(event.ReplaceCurrentSession{Path: i.Path()}),
			c.SendMsg(event.PopStackRoot{}),
		)
	}

	return m, nil
}

func (m *State) SaveSessionInputSubmit() (tea.Model, tea.Cmd) {
	v := strings.TrimSpace(m.saveInput.Value())
	if v == "" {
		return m, c.SendMsg(event.Err(errors.New("new name is empty")))
	}

	switch m.state {
	case val.STATE_SAVE_SESSION_INPUT:
		return m, tea.Batch(
			c.SendMsg(
				event.SaveSession{Path: fmt.Sprintf("%s/%s.json", val.CollectionFolder, v)},
			),
			c.SendMsg(event.PopStackRoot{}),
		)
	case val.STATE_SESSION_RENAME_INPUT:
		i, ok := m.sessionList.SelectedItem().(*store.Session)
		if !ok {
			return m, c.SendMsg(event.Err(errors.New("no session selected")))
		}

		err := i.Rename(v)
		if err != nil {
			return m, c.SendMsg(event.Err(err))
		}

		return m, tea.Batch(
			c.SendMsg(event.PopStack{}),
			c.SendMsg(event.LoadSessionList{}),
		)
	}

	return m, nil
}
