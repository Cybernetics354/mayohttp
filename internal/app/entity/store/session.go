package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Cybernetics354/mayohttp/internal/app/entity/intf"
	"github.com/Cybernetics354/mayohttp/internal/app/val"
)

func OpenSessionFromPath(path string) (*Session, error) {
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

type Session struct {
	Name          string              `json:"name"`
	Url           string              `json:"url"`
	Pipe          string              `json:"pipe"`
	PipedResponse string              `json:"piped_response"`
	Method        string              `json:"method"`
	Response      string              `json:"response"`
	Header        string              `json:"header"`
	Body          string              `json:"body"`
	ResFilter     intf.ResponseFilter `json:"res_filter"`
}

func (s Session) Title() string {
	return s.Name[:len(s.Name)-5]
}

func (s Session) Description() string {
	return fmt.Sprintf("%s: %s", s.Method, s.Url)
}

func (s Session) FilterValue() string {
	return s.Name + " " + s.Method + " " + s.Url + " " + s.Pipe
}

func (s *Session) Path() string {
	return fmt.Sprintf("%s/%s", val.CollectionFolder, s.Name)
}

func (s *Session) Delete() error {
	err := os.Remove(s.Path())
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) Rename(name string) error {
	err := os.Rename(s.Path(), fmt.Sprintf("%s/%s.json", val.CollectionFolder, name))
	if err != nil {
		return err
	}

	s.Name = name

	return nil
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
