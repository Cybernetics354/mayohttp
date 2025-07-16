package app

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Session struct {
	Name          string         `json:"name"`
	Url           string         `json:"url"`
	Pipe          string         `json:"pipe"`
	PipedResponse string         `json:"piped_response"`
	Method        string         `json:"method"`
	Response      string         `json:"response"`
	Header        string         `json:"header"`
	Body          string         `json:"body"`
	ResFilter     ResponseFilter `json:"res_filter"`
}

func createSessionFromState(s *State) *Session {
	return &Session{
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

func openSessionFromPath(path string) (*Session, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var session Session
	err = json.Unmarshal(b, &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *Session) Save(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

func (s *Session) Apply(m *State) *State {
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
