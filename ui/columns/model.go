package columns

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/louislouislouislouis/repr8ducer/k8s"
	"github.com/louislouislouislouis/repr8ducer/ui/commands"
	"github.com/louislouislouislouis/repr8ducer/ui/message"
)

const numberOfColums = 3

type ColumnsModel struct {
	IsActive        bool
	columns         []Column
	currentFocusCol colType
	k8sService      *k8s.K8sService
	generator       *k8s.Generator
}

type colType int

const (
	namespaceCol colType = iota
	podCol
	containerCol
)

func getListTitle(col colType) string {
	switch col {
	case namespaceCol:
		return "Namespaces"
	case podCol:
		return "Pods"
	default:
		return "Containers"
	}
}

func NewColumnsModel(namespace, pod, container string, generator *k8s.Generator) ColumnsModel {
	idxFocusedColumns := 0
	columns := make([]Column, numberOfColums)
	for idx, s := range []string{namespace, pod, container} {
		isColFocused := false
		if s == "" && idxFocusedColumns == 0 {
			idxFocusedColumns = idx
			isColFocused = true
		}
		columns[idx] = NewColumn(0, 0, isColFocused, false, setupCustomList(getListTitle(colType(idx)), []list.Item{}), s, spinner.New())
	}
	columns[idxFocusedColumns].SetFocus(true)
	return ColumnsModel{
		columns:         columns,
		currentFocusCol: colType(idxFocusedColumns),
		IsActive:        true,
		generator:       generator,
	}
}

func (m ColumnsModel) Init() tea.Cmd {
	var initCmd []tea.Cmd
	if m.columns[namespaceCol].current != "" {
		initCmd = append(
			initCmd,
			commands.GetPodsCmd(m.columns[namespaceCol].current, m.columns[podCol].current, context.TODO()),
		)
		if m.columns[podCol].current != "" {
			initCmd = append(
				initCmd,
				commands.GetContainersCmd(
					m.columns[namespaceCol].current,
					m.columns[podCol].current,
					m.columns[containerCol].current,
					context.TODO(),
				),
			)
		}
	}

	initCmd = append(
		initCmd,
		commands.GetNamespacesCmd(m.columns[namespaceCol].current, context.TODO()), // always init Namespace
		m.columns[podCol].spinner.Tick,
		m.columns[containerCol].spinner.Tick,
		m.columns[namespaceCol].spinner.Tick,
	)

	return tea.Batch(
		initCmd...,
	)
}

func (m ColumnsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.columns))
	switch msg := msg.(type) {

	case message.NamespaceMsg, message.PodMsg, message.ContainerMsg:
		return handlek8sMsg(m, msg)

	case tea.KeyMsg:
		if m.isFiltering() {
			break
		}
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			return m.handleEnterKey()

		case "left":
			return handleLeftKey(m)

		case "right":
			return m.handleEnterKey()
		}

	case tea.WindowSizeMsg:
		for i := range m.columns {
			col, cmd := m.columns[i].Update(msg)
			m.columns[i] = col
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case spinner.TickMsg:
		for i := range m.columns {
			col, cmd := m.columns[i].Update(msg)
			m.columns[i] = col
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	m.columns[m.currentFocusCol], cmd = m.columns[m.currentFocusCol].Update(msg)
	return m, tea.Batch(cmd)
}

func (m ColumnsModel) View() string {
	renders := make([]string, len(m.columns))
	for _, c := range m.columns {
		renders = append(renders, c.View())
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Left, renders...),
		m.columns[m.currentFocusCol].list.Help.View(m.columns[m.currentFocusCol].list),
	)
}

func setupCustomList(title string, items []list.Item) list.Model {
	theList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	theList.StatusMessageLifetime = time.Hour
	theList.SetShowHelp(false)
	theList.Title = title
	return theList
}

func (m ColumnsModel) isFiltering() bool {
	for i := range m.columns {
		if m.columns[i].list.FilterState() == list.Filtering {
			return true
		}
	}
	return false
}
