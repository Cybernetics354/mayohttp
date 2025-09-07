package ui

import (
	"slices"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type CompositeLayerPosition int

const (
	CompositeLayerTop CompositeLayerPosition = iota
	CompositeLayerBottom
	CompositeLayerLeft
	CompositeLayerRight
	CompositeLayerCenter
)

type CompositeViewLayer struct {
	Id        string
	View      string
	PositionX CompositeLayerPosition
	PositionY CompositeLayerPosition
	OffsetX   int
	OffsetY   int
}

func NewCompositeViewLayer() *CompositeViewLayer {
	return &CompositeViewLayer{
		Id:        "",
		View:      "",
		PositionX: CompositeLayerCenter,
		PositionY: CompositeLayerCenter,
		OffsetX:   0,
		OffsetY:   0,
	}
}

func (c *CompositeViewLayer) SetId(id string) *CompositeViewLayer {
	c.Id = id
	return c
}

func (c *CompositeViewLayer) SetView(view string) *CompositeViewLayer {
	c.View = view
	return c
}

func (c *CompositeViewLayer) SetPositionX(positionX CompositeLayerPosition) *CompositeViewLayer {
	c.PositionX = positionX
	return c
}

func (c *CompositeViewLayer) SetPositionY(positionY CompositeLayerPosition) *CompositeViewLayer {
	c.PositionY = positionY
	return c
}

func (c *CompositeViewLayer) SetOffset(offsetX int, offsetY int) *CompositeViewLayer {
	c.OffsetX, c.OffsetY = offsetX, offsetY
	return c
}

func (c *CompositeViewLayer) offsets(baseWidth int, baseHeight int) (int, int) {
	x, y := 0, 0

	switch c.PositionX {
	case CompositeLayerCenter:
		halfBaseWidth := baseWidth / 2
		halfLayerWidth := lipgloss.Width(c.View) / 2
		x = halfBaseWidth - halfLayerWidth
	case CompositeLayerRight:
		x = baseWidth - lipgloss.Width(c.View)
	}

	switch c.PositionY {
	case CompositeLayerCenter:
		halfBaseHeight := baseHeight / 2
		halfLayerHeight := lipgloss.Height(c.View) / 2
		y = halfBaseHeight - halfLayerHeight
	case CompositeLayerBottom:
		y = baseHeight - lipgloss.Height(c.View)
	}

	return x + c.OffsetX, y + c.OffsetY
}

type CompositeView struct {
	base   string
	layers []*CompositeViewLayer
}

func NewCompositeView(base string) *CompositeView {
	return &CompositeView{
		base: base,
	}
}

func (c *CompositeView) SetBase(base string) {
	c.base = base
}

func (c *CompositeView) AddLayer(layer *CompositeViewLayer) {
	c.layers = append(c.layers, layer)
}

func (c *CompositeView) RemoveLayer(id string) {
	if len(c.layers) == 0 {
		return
	}

	for i, layer := range c.layers {
		if layer.Id != id {
			continue
		}

		c.layers = slices.Delete(c.layers, i, i+1)
		break
	}
}

func (c *CompositeView) Render() string {
	res := c.base

	if len(c.layers) == 0 {
		return res
	}

	baseW, baseH := lipgloss.Size(res)

	for _, layer := range c.layers {
		w, h := lipgloss.Size(layer.View)
		if w >= baseW && h >= baseH {
			res = layer.View
			continue
		}

		x, y := layer.offsets(baseW, baseH)

		// clamp the layer so it doesn't go out of bounds
		x = clamp(x, 0, baseW-w)
		y = clamp(y, 0, baseH-h)

		baseLines := lines(res)
		layerLines := lines(layer.View)

		var sb strings.Builder
		for i, bgLine := range baseLines {
			//  add new line
			if i > 0 {
				sb.WriteString("\n")
			}

			// skip the replacements if the line is below or above the offset
			if i < y || i >= y+h {
				sb.WriteString(bgLine)
				continue
			}

			pos := 0
			if x > 0 {
				left := ansi.Truncate(bgLine, x, "")
				pos = ansi.StringWidth(left)
				sb.WriteString(left)
				if pos < x {
					sb.WriteString(whitespace(x - pos))
					pos = x
				}
			}

			fgLine := layerLines[i-y]
			sb.WriteString(fgLine)
			pos += ansi.StringWidth(fgLine)

			right := ansi.TruncateLeft(bgLine, pos, "")
			bgWidth := ansi.StringWidth(bgLine)
			rightWidth := ansi.StringWidth(right)
			if rightWidth <= bgWidth-pos {
				sb.WriteString(whitespace(bgWidth - rightWidth - pos))
			}
			sb.WriteString(right)
		}

		res = sb.String()
	}

	return res
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}

// Replace non-standard line endings with standard ones
func lines(view string) []string {
	s := strings.ReplaceAll(view, "\r\n", "\n")
	return strings.Split(s, "\n")
}

func whitespace(length int) string {
	return strings.Repeat(" ", length)
}
