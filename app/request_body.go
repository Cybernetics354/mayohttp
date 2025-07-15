package app

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type requestBody struct {
	raw  string
	form *multipart.Writer
}

func (r *requestBody) Buffer() (*bytes.Buffer, error) {
	reqType := r.ReqType()

	if reqType == REQUEST_BODY_FORM_SYMBOL {
		return r.FormBuffer()
	}

	return bytes.NewBuffer([]byte(r.Sanitized())), nil
}

func (r *requestBody) ReqType() string {
	if len(r.raw) == 0 {
		return REQUEST_BODY_RAW_SYMBOL
	}

	return strings.TrimSpace(r.raw[0:strings.Index(r.raw, "\n")])
}

func (r *requestBody) Sanitized() string {
	if reqType := r.ReqType(); reqType == REQUEST_BODY_FORM_SYMBOL ||
		reqType == REQUEST_BODY_RAW_SYMBOL {
		return strings.TrimSpace(r.raw[strings.Index(r.raw, "\n")+1:])
	}

	return r.raw
}

func (r *requestBody) FormBuffer() (*bytes.Buffer, error) {
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
			form.WriteField(key, val)
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
	r.form = form

	return &b, nil
}
