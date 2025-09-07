package intf

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cybernetics354/mayohttp/internal/app/val"
)

type RequestBody struct {
	Raw  string
	Form *multipart.Writer
}

func (r *RequestBody) Buffer() (*bytes.Buffer, error) {
	reqType := r.ReqType()

	if reqType == val.REQUEST_BODY_FORM_SYMBOL {
		return r.FormBuffer()
	}

	return bytes.NewBuffer([]byte(r.Sanitized())), nil
}

func (r *RequestBody) ReqType() string {
	if len(r.Raw) == 0 {
		return val.REQUEST_BODY_RAW_SYMBOL
	}

	return strings.TrimSpace(r.Raw[0:strings.Index(r.Raw, "\n")])
}

func (r *RequestBody) Sanitized() string {
	if reqType := r.ReqType(); reqType == val.REQUEST_BODY_FORM_SYMBOL ||
		reqType == val.REQUEST_BODY_RAW_SYMBOL {
		return strings.TrimSpace(r.Raw[strings.Index(r.Raw, "\n")+1:])
	}

	return r.Raw
}

func (r *RequestBody) FormBuffer() (*bytes.Buffer, error) {
	rawForms := strings.Split(r.Sanitized(), "\n")

	var b bytes.Buffer

	form := multipart.NewWriter(&b)

	for _, line := range rawForms {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		keyValLine := strings.SplitN(trimmed, "=", 2)
		if len(keyValLine) == 2 {
			key, val := keyValLine[0], keyValLine[1]
			err := form.WriteField(key, val)
			if err != nil {
				return nil, err
			}

			continue
		}

		fileLine := strings.SplitN(trimmed, "@", 2)
		if len(fileLine) == 2 {
			name, path := fileLine[0], fileLine[1]

			file, err := os.Open(path)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			f, err := form.CreateFormFile(name, filepath.Base(path))
			if err != nil {
				return nil, err
			}

			if _, err := io.Copy(f, file); err != nil {
				return nil, err
			}

			continue
		}
	}

	form.Close()
	r.Form = form

	return &b, nil
}
