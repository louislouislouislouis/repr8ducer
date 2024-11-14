package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	v1 "k8s.io/api/core/v1"

	"github.com/louislouislouislouis/repr8ducer/k8s"
)

type model struct {
	k8sService    *k8s.K8sService
	namespace     string
	pod           string
	container     string
	listNamespace list.Model
	listContainer list.Model
	listPods      list.Model
}

func NewModel(k8s *k8s.K8sService, namespace, podName, container string) model {
	displayedNamespaceList := []list.Item{}
	displayedPodList := []list.Item{}
	displayedContainerList := []list.Item{}

	fmt.Println(fmt.Sprintln(namespace == "", podName == "", container == ""))
	if namespace == "" {
		listNamespace, _ := k8s.ListNamespace()
		displayedNamespaceList = createDisplayedListFromMetadata(
			listNamespace.Items,
			func(nms v1.Namespace) Descripable {
				return &descipableMeta{&nms.ObjectMeta}
			},
		)
	}

	fmt.Println(fmt.Sprintln(namespace == "", podName == "", container == ""))
	if podName == "" && namespace != "" {
		pods, err := k8s.ListPodsInNamespace(namespace)
		if err != nil {
			panic(err.Error())
		}
		displayedPodList = createDisplayedListFromMetadata(
			pods.Items,
			func(nms v1.Pod) Descripable {
				return &descipableMeta{&nms.ObjectMeta}
			},
		)
	}

	fmt.Println(fmt.Sprintln(namespace == "", podName == "", container == ""))
	if podName == "p" && namespace == "" && container == "" {

		container, _ := k8s.GetContainerFromPods(podName, namespace)
		displayedContainerList = createDisplayedListFromMetadata(
			container,
			func(container v1.Container) Descripable {
				return &descipableContainer{container}
			},
		)
	}

	return model{
		k8sService:    k8s,
		pod:           podName,
		namespace:     namespace,
		container:     container,
		listNamespace: list.New(displayedNamespaceList, list.NewDefaultDelegate(), 0, 0),
		listPods:      list.New(displayedPodList, list.NewDefaultDelegate(), 0, 0),
		listContainer: list.New(displayedContainerList, list.NewDefaultDelegate(), 0, 0),
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

	if m.container == "" {
		return docStyle.Render(m.listContainer.View())
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
		m.listContainer.SetSize(msg.Width-h, msg.Height-v)
		m.listPods.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmdNamespace, cmdPods, cmdContainer tea.Cmd
	m.listNamespace, cmdNamespace = m.listNamespace.Update(msg)
	m.listPods, cmdPods = m.listPods.Update(msg)
	m.listContainer, cmdContainer = m.listContainer.Update(msg)
	return m, tea.Batch(cmdPods, cmdNamespace, cmdContainer)
}

func (m model) onNamespaceSelection() (tea.Model, tea.Cmd) {
	m.namespace = m.listNamespace.Items()[m.listNamespace.Index()].(DisplayedItem).title
	pods, err := m.k8sService.ListPodsInNamespace(m.namespace)
	if err != nil {
		panic(err.Error())
	}
	cmd := m.listPods.SetItems(
		createDisplayedListFromMetadata(pods.Items, func(pod v1.Pod) Descripable {
			return &descipableMeta{&pod}
		}),
	)
	return m, cmd
}

func (m model) onPodSelection() (tea.Model, tea.Cmd) {
	m.pod = m.listPods.Items()[m.listPods.Index()].(DisplayedItem).title
	// strcmd, _ := m.k8sService.Exec(m.namespace, m.pod)
	container, _ := m.k8sService.GetContainerFromPods(m.pod, m.namespace)
	displayedContainerList := createDisplayedListFromMetadata(
		container,
		func(c v1.Container) Descripable {
			return &descipableContainer{c}
		},
	)
	cmd := m.listContainer.SetItems(displayedContainerList)

	return m, cmd
}
