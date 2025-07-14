package app

import (
	"fmt"
	"strings"

	"github.com/Cybernetics354/mayohttp/app/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusResponseFilter = lipgloss.NewStyle().Foreground(ui.FocusColor)
	blurResponseFilter  = lipgloss.NewStyle().Foreground(ui.BlurColor)
)

type ResponseFilter struct {
	ReqHeader     bool `json:"req_header,omitempty"`
	ReqBody       bool `json:"req_body,omitempty"`
	ResHeader     bool `json:"res_header,omitempty"`
	ResBody       bool `json:"res_body,omitempty"`
	PositionIndex int  `json:"position_index,omitempty"`
	focus         bool
}

func CreateResponseFilter() ResponseFilter {
	return ResponseFilter{
		ReqHeader:     true,
		ReqBody:       true,
		ResHeader:     true,
		ResBody:       true,
		focus:         false,
		PositionIndex: 0,
	}
}

func (r *ResponseFilter) Length() int {
	return 4
}

func (r *ResponseFilter) Next() {
	r.PositionIndex = min(r.PositionIndex+1, r.Length()-1)
}

func (r *ResponseFilter) Prev() {
	r.PositionIndex = max(r.PositionIndex-1, 0)
}

func (r *ResponseFilter) Toggle() {
	switch r.PositionIndex {
	case 0:
		r.ReqHeader = !r.ReqHeader
	case 1:
		r.ReqBody = !r.ReqBody
	case 2:
		r.ResHeader = !r.ResHeader
	case 3:
		r.ResBody = !r.ResBody
	}
}

func (r *ResponseFilter) renderField(str string, active bool, focus bool) string {
	point := "[ ]"
	if active {
		point = "[*]"
	}

	if r.focus && focus {
		return focusResponseFilter.Render(fmt.Sprintf("%s %s", point, str))
	}

	return blurResponseFilter.Render(fmt.Sprintf("%s %s", point, str))
}

func (r *ResponseFilter) Focus() {
	r.focus = true
}

func (r *ResponseFilter) Blur() {
	r.focus = false
}

func (r *ResponseFilter) Render() string {
	strList := []string{}

	strList = append(strList, r.renderField("Req Header", r.ReqHeader, r.PositionIndex == 0))
	strList = append(strList, r.renderField("Req Body", r.ReqBody, r.PositionIndex == 1))
	strList = append(strList, r.renderField("Res Header", r.ResHeader, r.PositionIndex == 2))
	strList = append(strList, r.renderField("Res Body", r.ResBody, r.PositionIndex == 3))

	return strings.Join(strList, " • ")
}

func (r *ResponseFilter) HandleKeyPress(msg tea.KeyMsg) (ResponseFilter, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "h", "left":
		r.Prev()
	case "l", "right":
		r.Next()
	case " ":
		r.Toggle()
		cmd = sendMsg(runPipeMsg{})
	}

	return *r, cmd
}

func (r *ResponseFilter) Filter(str string) string {
	list := strings.Split(str, responseSeparator)
	if len(list) != r.Length() {
		return str
	}

	for i := range list {
		switch i {
		case 0:
			if !r.ReqHeader {
				list[i] = ""
			}
		case 1:
			if !r.ReqBody {
				list[i] = ""
			}
		case 2:
			if !r.ResHeader {
				list[i] = ""
			}
		case 3:
			if !r.ResBody {
				list[i] = ""
			}
		}
	}

	// remove empty strings
	var filtered []string
	for _, v := range list {
		if v != "" {
			filtered = append(filtered, v)
		}
	}

	return strings.Join(filtered, responseSeparator)
}
