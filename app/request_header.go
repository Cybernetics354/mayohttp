package app

import (
	"net/http"
	"strings"
)

type requestHeader struct {
	raw string
}

func (r *requestHeader) Apply(req *http.Request) {
	lines := strings.Split(r.raw, "\n")

	for _, line := range lines {
		header := strings.SplitN(line, ":", 2)
		if len(header) != 2 {
			continue
		}

		key, val := strings.TrimSpace(header[0]), strings.TrimSpace(header[1])
		req.Header.Add(key, val)
	}
}
