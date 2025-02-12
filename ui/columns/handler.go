package columns

import (
	"context"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/louislouislouislouis/repr8ducer/ui/commands"
	"github.com/louislouislouislouis/repr8ducer/ui/common"
	"github.com/louislouislouislouis/repr8ducer/ui/message"
	"github.com/louislouislouislouis/repr8ducer/ui/modifier"
	statusline "github.com/louislouislouislouis/repr8ducer/ui/statusLine"
	"github.com/louislouislouislouis/repr8ducer/utils"
	v1 "k8s.io/api/core/v1"
)

type direction int

const (
	left direction = iota
	right
)

func (m ColumnsModel) changeFocus(direction direction) ColumnsModel {
	if direction == left {
		m.currentFocusCol--
	} else {
		m.currentFocusCol++
	}

	if m.currentFocusCol < 0 {
		m.currentFocusCol = 0
	}
	if m.currentFocusCol > containerCol {
		m.currentFocusCol = containerCol
	}

	for i := range m.columns {
		m.columns[i].isFocused = i == int(m.currentFocusCol)
	}

	return m
}

func handleLeftKey(m ColumnsModel) (tea.Model, tea.Cmd) {
	return m.changeFocus(left), nil
}

func handlek8sMsg(m ColumnsModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case message.NamespaceMsg:
		items := common.CreateDisplayedListFromMetadata(msg.Val.List.Items, func(nms v1.Namespace) common.DiplayableItemList {
			return &common.DisplayableMeta{&nms.ObjectMeta}
		})
		m.columns[namespaceCol], cmd = m.columns[namespaceCol].Update(
			message.ListUpdateMsg{
				StatusTxt:        "namespace",
				Title:            "namespace",
				Val:              items,
				PreSelectedValue: msg.Val.PreSelectedNamespace,
			},
		)

	case message.PodMsg:
		items := common.CreateDisplayedListFromMetadata(msg.Val.List.Items, func(nms v1.Pod) common.DiplayableItemList {
			return &common.DisplayableMeta{&nms.ObjectMeta}
		})
		m.columns[podCol], cmd = m.columns[podCol].Update(
			message.ListUpdateMsg{
				StatusTxt:        "pod",
				Title:            "Pod",
				Val:              items,
				PreSelectedValue: msg.Val.PreSelectedPod,
			},
		)

	case message.ContainerMsg:
		items := common.CreateDisplayedListFromMetadata(msg.Val.List, func(container v1.Container) common.DiplayableItemList {
			return &common.DisplayableContainer{container}
		})
		m.columns[containerCol], cmd = m.columns[containerCol].Update(
			message.ListUpdateMsg{
				StatusTxt:        "hehe",
				Title:            "Container",
				Val:              items,
				PreSelectedValue: msg.Val.PreSelectedContainer,
			},
		)
	}

	return m, cmd
}

func (m ColumnsModel) handleEnterKey() (ColumnsModel, tea.Cmd) {
	var cmd tea.Cmd

	if len(m.columns[m.currentFocusCol].list.VisibleItems()) == 0 {
		return m, nil
	}

	m.columns[m.currentFocusCol].current = selectTitleSelected(m.columns[m.currentFocusCol].list)
	switch m.currentFocusCol {
	case podCol:
		cmd = commands.GetContainersCmd(
			m.columns[namespaceCol].current,
			m.columns[podCol].current,
			"",
			context.TODO(),
		)
	case containerCol:
		command, err := m.generator.PodToContainer(
			m.columns[namespaceCol].current,
			m.columns[podCol].current,
			context.TODO(),
		)
		if err != nil {
			utils.Log.Error().Msg(err.Error())
			cmd = statusline.UpdateStatusLine(err.Error())
			return m, cmd
		}
		clipboard.WriteAll(command.GetCommand())
		if len(command.Modifiers) != 0 {
			// TODO Open dialog
			cmds := []tea.Cmd{}
			for _, cmd := range command.Modifiers {
				keys := make([]string, len(cmd.GetDetections()))
				i := 0
				for k := range cmd.GetDetections() {
					keys[i] = k
					i++
				}
				cmds = append(cmds, modifier.UpdateModifierRecap(keys), commands.ChangeMainModelFocus(1))
			}
			cmd = tea.Batch(cmds...)
		}

	case namespaceCol:
		cmd = commands.GetPodsCmd(m.columns[namespaceCol].current, "", context.TODO())

	}
	m = m.changeFocus(right)
	return m, cmd
}

func selectTitleSelected(list list.Model) string {
	return list.VisibleItems()[list.Index()].(common.DisplayedItem).Title()
}
