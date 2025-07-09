package app

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func formatResponse(req *http.Request) string {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}

	body := fmt.Sprintf("RESPONSE:\n%s", string(respDump))

	return body
}
