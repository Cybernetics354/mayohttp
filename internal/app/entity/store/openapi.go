package store

import (
	"errors"
	"os"

	"github.com/pb33f/libopenapi"
)

type OpenAPI struct {
	Path string
	doc  libopenapi.Document
}

func (o *OpenAPI) Open() error {
	file, err := os.ReadFile(o.Path)
	if err != nil {
		return err
	}

	doc, err := libopenapi.NewDocument(file)
	if err != nil {
		return err
	}

	o.doc = doc

	return nil
}

func (o *OpenAPI) GetDoc() (libopenapi.Document, error) {
	if o.doc == nil {
		return nil, errors.New("Please run the .Open first")
	}

	return o.doc, nil
}
