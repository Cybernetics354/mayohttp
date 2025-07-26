package ui

func RenderWithHeader(base string, header string) string {
	view := NewCompositeView(base)
	headerLayer := NewCompositeViewLayer().
		SetPositionY(CompositeLayerTop).
		SetPositionX(CompositeLayerLeft).
		SetOffset(2, 0).
		SetView(header)
	view.AddLayer(headerLayer)

	return view.Render()
}
