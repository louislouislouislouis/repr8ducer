package mainmodel

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/louislouislouislouis/repr8ducer/k8s"
	"github.com/louislouislouislouis/repr8ducer/ui/columns"
	"github.com/louislouislouislouis/repr8ducer/ui/message"
	"github.com/louislouislouislouis/repr8ducer/ui/modifier"
	statusline "github.com/louislouislouislouis/repr8ducer/ui/statusLine"
	"github.com/louislouislouislouis/repr8ducer/ui/styles"
)

type itemWithList struct {
	current       string
	list          list.Model
	isInitialized bool
}

type mainModel struct {
	width            int
	colView          columns.ColumnsModel
	modifierView     modifier.ModifierModel
	statusline       statusline.StatusLine
	k8sService       *k8s.K8sService
	generator        *k8s.Generator
	SkipUrlOverwrite bool
}

type MainModelConfig struct {
	Namespace, Pod, Container string
	SkipUrlOverwrite          bool
}

type MainModelMode int

const (
	colView MainModelMode = iota
	modifierView
)

func NewMainModel(k8sService *k8s.K8sService, c MainModelConfig) mainModel {
	generator := k8s.NewDefaultGenerator(k8sService)
	return mainModel{
		colView:          columns.NewColumnsModel(c.Namespace, c.Pod, c.Container, generator),
		k8sService:       k8sService,
		generator:        generator,
		modifierView:     modifier.ModifierModel{},
		statusline:       statusline.StatusLine{Text: "Everything is good until now"},
		SkipUrlOverwrite: c.SkipUrlOverwrite,
	}
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(
		m.colView.Init(),
		m.modifierView.Init(),
		m.statusline.Init(),
	)
}

func (m mainModel) View() string {
	activeModel, _ := m.getFocusPart()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title(m.width),
		activeModel.View(),
		m.statusline.View(),
	)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	lastActiveModel, mode := m.getFocusPart()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		// No return because child model also need to be set
	case message.InfoMsg:
		statusline, cmd := m.statusline.Update(msg)
		m.statusline = statusline
		return m, cmd
	case message.MainModelChangeFocusMsg:
		m = m.setFocus(MainModelMode(msg.Val))
		if m.SkipUrlOverwrite {
			if _, mode := m.getFocusPart(); mode == modifierView {
				return m, tea.Quit
			}
		}
		return m, nil
	}
	switch mode {
	case colView:
		if lastColView, ok := lastActiveModel.(columns.ColumnsModel); ok {
			newActiveModel, cmd := lastColView.Update(msg)
			m.colView = newActiveModel.(columns.ColumnsModel)
			return m, cmd
		}
	case modifierView:
		if lastModifierView, ok := lastActiveModel.(modifier.ModifierModel); ok {
			newActiveModel, cmd := lastModifierView.Update(msg)
			m.modifierView = newActiveModel.(modifier.ModifierModel)
			return m, cmd
		}
	}
	return m, nil
}

func (m mainModel) setFocus(mode MainModelMode) mainModel {
	m.colView.IsActive = mode == colView
	m.modifierView.IsActive = mode == modifierView
	return m
}

func (m mainModel) getFocusPart() (tea.Model, MainModelMode) {
	// TODO CHANGE THIS
	if m.colView.IsActive {
		return m.colView, colView
	}
	return m.modifierView, modifierView
}

func title(width int) string {
	title := `  ____  ___  ___  ____  
 (_  _)(  _)/ __)(_  _) 
   )(   ) _)\__ \  )(   
(__) (___)(___/ (__)    `
	return styles.BigTitleStyle.Width(width).Render(title)
}
