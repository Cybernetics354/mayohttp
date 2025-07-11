package app

type Session struct {
	Url           string `json:"url"`
	Pipe          string `json:"pipe"`
	PipedResponse string `json:"piped_response"`
	Method        string `json:"method"`
	Response      string `json:"response"`
	Header        string `json:"header"`
	Body          string `json:"body"`
}
