package app

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
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

func getDefaultEditor() string {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	return editor
}
