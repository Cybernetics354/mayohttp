package intf

type FileItem struct {
	Name string
	Path string
}

func (f FileItem) Title() string {
	return f.Name
}

func (f FileItem) Description() string {
	return f.Path
}

func (f FileItem) FilterValue() string {
	return f.Path
}
