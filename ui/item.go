package ui

type DisplayedItem struct {
	title, desc string
	isSelected  bool
}

func (i DisplayedItem) Title() string       { return i.title }
func (i DisplayedItem) Description() string { return i.desc }
func (i DisplayedItem) FilterValue() string { return i.title }
func (i DisplayedItem) IsSelected() bool    { return i.isSelected }
