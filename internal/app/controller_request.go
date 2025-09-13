package app

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	c "github.com/Cybernetics354/mayohttp/internal/app/common"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/event"
	"github.com/Cybernetics354/mayohttp/internal/app/entity/intf"
	"github.com/Cybernetics354/mayohttp/internal/app/val"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *State) HandleRequestResult(msg event.RequestResult) (tea.Model, tea.Cmd) {
	m.response.SetValue(strings.TrimSpace(msg.Res))
	return m, tea.Batch(
		c.SendMsg(event.HideSpinner{}),
		c.SendMsg(event.SetActivity("Request complete")),
		c.SendMsg(event.RunPipe{}),
		c.SendMsg(event.SaveSession{Path: val.DefaultSessionPath}),
		event.ListenResponseCmd(m.resSub),
	)
}

func (m *State) HandlePipeResult(msg event.PipeResult) (tea.Model, tea.Cmd) {
	m.pipedresp.SetValue(strings.TrimSpace(msg.Res))
	return m, tea.Batch(
		c.SendMsg(event.HideSpinner{}),
		c.SendMsg(event.SetActivity("Piping complete")),
		c.SendMsg(event.SaveSession{Path: val.DefaultSessionPath}),
		event.ListenPipeResponseCmd(m.pipeResSub),
	)
}

func (m *State) RunRequest() (tea.Model, tea.Cmd) {
	return m, tea.Batch(
		c.SendMsg(event.ShowSpinner{}),
		c.SendMsg(event.SetActivity("Requesting...")),
		m.Request,
	)
}

func (m *State) Request() tea.Msg {
	uc := make(chan c.ParseWithEnvMsg)
	bc := make(chan c.ParseWithEnvMsg)
	hc := make(chan c.ParseWithEnvMsg)
	defer close(uc)
	defer close(bc)
	defer close(hc)

	go c.ParseWithEnv(m.url.Value(), uc)
	go c.ParseWithEnv(m.body.Value(), bc)
	go c.ParseWithEnv(m.header.Value(), hc)

	url := <-uc
	if url.Err != nil {
		m.resSub <- event.RequestResult{
			Err: url.Err,
			Res: fmt.Sprintf("URL parse error : %s", url.Err.Error()),
		}
		return nil
	}

	body := <-bc
	if body.Err != nil {
		m.resSub <- event.RequestResult{
			Err: body.Err,
			Res: fmt.Sprintf("Body parse error : %s", body.Err.Error()),
		}
		return nil
	}

	header := <-hc
	if header.Err != nil {
		m.resSub <- event.RequestResult{
			Err: header.Err,
			Res: fmt.Sprintf("Header parse error : %s", header.Err.Error()),
		}
		return nil
	}

	reqBody := intf.RequestBody{Raw: body.Str}
	bodyReader, err := reqBody.Buffer()
	if err != nil {
		m.resSub <- event.RequestResult{
			Err: err,
			Res: fmt.Sprintf("Invalid body : %s", err.Error()),
		}
		return nil
	}

	if reqBody.Form != nil {
		header.Str = fmt.Sprintf(
			"Content-Type: %s\n%s",
			reqBody.Form.FormDataContentType(),
			header.Str,
		)
	}

	req, err := http.NewRequest(m.method, url.Str, bodyReader)
	if err != nil {
		m.resSub <- event.RequestResult{
			Err: err,
			Res: fmt.Sprintf("Request error : %s", err.Error()),
		}
		return nil
	}

	// The header example will look like this:
	// Header-1: value1
	// Header-2: value2
	reqHeader := intf.RequestHeader{Raw: header.Str}
	reqHeader.Apply(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		m.resSub <- event.RequestResult{
			Err: err,
			Res: fmt.Sprintf("Request do error : %s", err.Error()),
		}
		return nil
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		m.resSub <- event.RequestResult{
			Err: err,
			Res: fmt.Sprintf("Response dump error : %s", err.Error()),
		}
		return nil
	}

	var resBuffer bytes.Buffer
	resBuffer.WriteString(header.Str + "\n")
	resBuffer.WriteString(val.ResponseSeparator + "\n")
	resBuffer.WriteString(body.Str + "\n")
	resBuffer.WriteString(val.ResponseSeparator + "\n")
	for k, v := range resp.Header {
		resBuffer.WriteString(fmt.Sprintf("%s: %s\n", k, v))
	}
	resBuffer.WriteString(val.ResponseSeparator + "\n")
	resBuffer.WriteString(string(bodyBytes) + "\n")

	m.resSub <- event.RequestResult{Res: resBuffer.String()}
	return nil
}

func (m *State) RunPipe() (tea.Model, tea.Cmd) {
	return m, tea.Batch(
		c.SendMsg(event.ShowSpinner{}),
		c.SendMsg(event.SetActivity("Piping...")),
		m.PipeRequest,
	)
}

func (m *State) PipeRequest() tea.Msg {
	resp, pipe := m.response.Value(), m.pipe.Value()
	if resp == "" {
		m.pipeResSub <- event.PipeResult{Res: resp}
		return nil
	}

	resp = m.resFilter.Filter(resp)
	if pipe == "" {
		m.pipeResSub <- event.PipeResult{Res: resp}
		return nil
	}

	command := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"set -a && source %s && set +a && echo '%s' | %s",
			val.EnvFilePath,
			strings.ReplaceAll(resp, "'", "'\\''"),
			pipe,
		),
	)
	output, err := command.CombinedOutput()
	if err != nil {
		m.pipeResSub <- event.PipeResult{Err: err, Res: string(output)}
		return nil
	}

	m.pipeResSub <- event.PipeResult{Res: string(output)}
	return nil
}
