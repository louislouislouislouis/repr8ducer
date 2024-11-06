package ui

import (
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/louislouislouislouis/repr8ducer/k8s"
)

type model struct {
	k8sService    *k8s.K8sService
	namespace     string
	pod           string
	listNamespace list.Model
	listPods      list.Model
}

func NewModel(k8s *k8s.K8sService) model {
	listeNamespace, _ := k8s.ListNamespace()
	liste := createDisplayedListFromMetadata(listeNamespace.Items, func(nms v1.Namespace) metaV1.Object {
		return &nms.ObjectMeta
	})

	return model{
		k8sService:    k8s,
		listNamespace: list.New(liste, list.NewDefaultDelegate(), 0, 0),
		listPods:      list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	if m.namespace == "" {
		return docStyle.Render(m.listNamespace.View())
	}

	if m.pod == "" {
		return docStyle.Render(m.listPods.View())
	}

	return "Namespace: " + m.namespace + " Pod: " + m.pod
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			if m.namespace == "" {
				return m.onNamespaceSelection()
			}
			if m.pod == "" {
				return m.onPodSelection()
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.listNamespace.SetSize(msg.Width-h, msg.Height-v)
		m.listPods.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmdNamespace, cmdPods tea.Cmd
	m.listNamespace, cmdNamespace = m.listNamespace.Update(msg)
	m.listPods, cmdPods = m.listPods.Update(msg)

	return m, tea.Batch(cmdPods, cmdNamespace)
}

func (m model) onNamespaceSelection() (tea.Model, tea.Cmd) {
	m.namespace = m.listNamespace.Items()[m.listNamespace.Index()].(DisplayedItem).title
	pods, err := m.k8sService.ListPodsInNamespace(m.namespace)
	if err != nil {
		panic(err.Error())
	}
	cmd := m.listPods.SetItems(createDisplayedListFromMetadata(pods.Items, func(pod v1.Pod) metaV1.Object {
		return &pod.ObjectMeta
	}))
	return m, cmd
}

func (m model) onPodSelection() (tea.Model, tea.Cmd) {
	m.pod = m.listPods.Items()[m.listPods.Index()].(DisplayedItem).title
	strcmd, _ := m.k8sService.Exec(m.namespace, m.pod)

	// Copy command to copy paste buffer
	// fmt.Println(strcmd)
	clipboard.WriteAll(strcmd)

	return m, tea.Quit
}
