package intf

import (
	"net/http"
	"strings"
)

type RequestHeader struct {
	Raw string
}

func (r *RequestHeader) Apply(req *http.Request) {
	lines := strings.SplitSeq(r.Raw, "\n")

	for line := range lines {
		header := strings.SplitN(line, ":", 2)
		if len(header) != 2 {
			continue
		}

		key, val := strings.TrimSpace(header[0]), strings.TrimSpace(header[1])
		req.Header.Add(key, val)
	}
}
