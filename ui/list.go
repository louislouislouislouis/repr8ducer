package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func NewModel(liste []list.Item) model {
	d := list.NewDefaultDelegate()
	d.ShowDescription = true

	return model{
		listNamespace: list.New(liste, d, 0, 0),
		listPods:      list.New([]list.Item{}, d, 0, 0),
	}
}

type model struct {
	namespace     string
	pod           string
	listNamespace list.Model
	listPods      list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

type customDelegate struct {
	list.DefaultDelegate
}

type item struct {
	title, desc string
	isSelected  bool
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
func (i item) IsSelected() bool    { return i.isSelected }

// View personnalise l'affichage de chaque élément dans la liste
func (d customDelegate) View(i list.Item) string {
	it := i.(item)
	if it.isSelected {
		return selectedStyle.Render(it.Title())
	}
	return unselectedStyle.Render(it.Title())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			fmt.Println("forh")
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			if m.namespace == "" {
				index := m.listNamespace.Index()
				item44 := m.listNamespace.Items()[index].(item)
				m.namespace = item44.title

				m.listPods.SetItems(CreateFakeList(44, "Pods bb"))
				return m, nil
			}

			index := m.listPods.Index()
			item44 := m.listPods.Items()[index].(item)
			m.pod = item44.title
			return m, nil

		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.listNamespace.SetSize(msg.Width-h, msg.Height-v)
		m.listPods.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.listNamespace, cmd = m.listNamespace.Update(msg)
	var cmd2 tea.Cmd
	m.listPods, cmd2 = m.listPods.Update(msg)

	return m, tea.Batch(cmd2, cmd)
}

func (m model) View() string {
	toretrun := "d\n\né\n"
	if m.namespace == "" {
		return docStyle.Render(toretrun + m.listNamespace.View())
	}
	toretrun += "Namespacdjend e:\n " + m.namespace

	if m.pod == "" {
		return toretrun + docStyle.Render(toretrun+m.listPods.View())
	}

	return "Namespace: " + m.namespace + " Pod: " + m.pod
}

func CreateFakeList(numItems int, descr string) []list.Item {
	items := make([]list.Item, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = item{
			title: fmt.Sprintf("%s %d", descr, i+1),
			desc:  fmt.Sprintf("Description for %s %d", descr, i+1),
		}
	}
	return items
}
