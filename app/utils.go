package app

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func formatResponse(req *http.Request) string {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Sprintf("Request error : %s", err.Error())
	}

	defer resp.Body.Close()

	respDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return fmt.Sprintf("Response dump error : %s", err.Error())
	}

	return string(respDump)
}
