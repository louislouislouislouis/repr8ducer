package commands

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/louislouislouislouis/repr8ducer/k8s"
	"github.com/louislouislouislouis/repr8ducer/modifiers"
	"github.com/louislouislouislouis/repr8ducer/ui/message"
)

func GetNamespacesCmd(preSelectedNamepace string, ctx context.Context) tea.Cmd {
	return func() tea.Msg {
		// Todo hadle error
		listNamespace, _ := k8s.GetService().ListNamespace(ctx)
		return message.NamespaceMsg{
			Val: message.NamespaceMsgValue{
				List:                 listNamespace,
				PreSelectedNamespace: preSelectedNamepace,
			},
		}
	}
}

func GetContainersCmd(nms, pod, preSelectedContainer string, ctx context.Context) tea.Cmd {
	// Todo hadle error
	return func() tea.Msg {
		containers, _ := k8s.GetService().GetContainerFromPods(nms, pod, ctx)
		return message.ContainerMsg{
			Val: message.ContainerMsgValue{
				List:                 containers,
				PreSelectedContainer: preSelectedContainer,
			},
		}
	}
}

func GetPodsCmd(nms, preSelectedPod string, ctx context.Context) tea.Cmd {
	// Todo hadle error
	return func() tea.Msg {
		listPods, _ := k8s.GetService().ListPodsInNamespace(nms, ctx)
		return message.PodMsg{
			Val: message.PodMsgValue{
				List:           listPods,
				PreSelectedPod: preSelectedPod,
			},
		}
	}
}

func ChangeMainModelFocus(mode int) tea.Cmd {
	return func() tea.Msg {
		return message.MainModelChangeFocusMsg{
			Val: mode,
		}
	}
}

func GetUrlsFromFolder(folderPath string) tea.Cmd {
	return func() tea.Msg {
		msgValue, err := modifiers.NewUrlReplacer().SearchUrlsInDir(folderPath)
		if err != nil {
			return message.UrlDetectionMsg{
				Status: message.Error,
			}
		}
		return message.UrlDetectionMsg{
			Status: message.Ok,
			Val:    msgValue,
		}
	}
}
