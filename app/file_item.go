package app

type fileItem struct {
	name string
	path string
}

func (f fileItem) Title() string {
	return f.name
}

func (f fileItem) Description() string {
	return f.path
}

func (f fileItem) FilterValue() string {
	return f.path
}
