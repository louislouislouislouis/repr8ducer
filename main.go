package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/louislouislouislouis/repr8ducer/k8s"
)

var DEFAULT_NAMESPACE = "things"

// func main() {
// 	k8sService := k8s.GetService()
// 	//k8sService.ListNamespace()
// 	//k8sService.ListPodsInNamespace(DEFAULT_NAMESPACE)
// 	err := k8sService.Exec()
// 	fmt.Println(err)
// }

var (
	docStyle   = lipgloss.NewStyle().Margin(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")). // Couleur blanche
			Background(lipgloss.Color("#F0F")).    // Fond bleu
			Bold(true).
			Padding(1, 2)
	selectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0")).Bold(true)
	unselectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
)

type item struct {
	title, desc string
	isSelected  bool
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
func (i item) IsSelected() bool    { return i.isSelected }

type model struct {
	list      list.Model
	textInput textinput.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

type customDelegate struct {
	list.DefaultDelegate
}

// View personnalise l'affichage de chaque élément dans la liste
func (d customDelegate) View(i list.Item) string {
	it := i.(item) // Assurez-vous que c'est le bon type
	if it.isSelected {
		return selectedStyle.Render(it.Title())
	}
	return unselectedStyle.Render(it.Title())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		// fmt.Println(msg.String())

		if msg.String() == "ctrl+c" {
			fmt.Println("forh")
			return m, tea.Quit
		}
		if msg.String() == " " {
			index := m.list.Index()
			item44 := m.list.Items()[index].(item)
			// os.Exit(0)
			// fmt.Println(item44.title
			item44.title = "zdsfzes"
			cmd := m.list.SetItem(index, item{
				title: item44.title,
				desc:  item44.desc,
			})
			// quitter l'appli tea
			cmd = tea.Quit

			return m, cmd
			// em44.isSelected = true
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	// fmt.Println("cc")
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	k8sService := k8s.GetService()
	itemss, _ := k8sService.ListNamespace()
	items := []list.Item{}

	for _, o := range itemss.Items {
		newItem := item{
			title: o.Name,             // Utilisez le champ que vous voulez pour le tqitre
			desc:  string(o.GetUID()), // Description personnalisée
		}
		items = append(items, newItem)
	}
	ti := textinput.New()
	ti.Placeholder = "Ajouter un nouvel élément"
	ti.Focus()
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	list := list.New(items, d, 20, 10)
	m := model{list: list, textInput: ti}
	m.list.Title = "My Fave Things  fzef"
	m.list.Styles.Title = titleStyle
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	fmt.Println("cc")
}
