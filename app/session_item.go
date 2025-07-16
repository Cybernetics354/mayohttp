package app

import "fmt"

type SessionItem struct {
	name    string
	session *Session
}

func (s SessionItem) Title() string {
	return s.name[:len(s.name)-5]
}

func (s SessionItem) Description() string {
	return fmt.Sprintf("%s: %s", s.session.Method, s.session.Url)
}

func (s SessionItem) FilterValue() string {
	return s.name + " " + s.session.Method + " " + s.session.Url + " " + s.session.Pipe
}

func (s *SessionItem) Path() string {
	return fmt.Sprintf("%s/%s", collectionFolder, s.name)
}

func (s *SessionItem) GetSession() error {
	session, err := openSessionFromPath(s.Path())
	if err != nil {
		return err
	}

	s.session = session

	return nil
}
