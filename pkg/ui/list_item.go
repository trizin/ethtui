package ui

type ListItem struct {
	title string
	desc  string
	id    string
}

func (i ListItem) Title() string       { return i.title }
func (i ListItem) Description() string { return i.desc }
func (i ListItem) FilterValue() string { return i.title }
