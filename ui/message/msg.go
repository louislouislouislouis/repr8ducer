package message

import (
	"github.com/charmbracelet/bubbles/list"
	v1 "k8s.io/api/core/v1"
)

type StatusMsg int

const (
	Error StatusMsg = iota
	Ok
)

type NamespaceMsgValue struct {
	List                 *v1.NamespaceList
	PreSelectedNamespace string
}

type ContainerMsgValue struct {
	List                 []v1.Container
	PreSelectedContainer string
}
type PodMsgValue struct {
	List           *v1.PodList
	PreSelectedPod string
}

type NamespaceMsg struct {
	Val    NamespaceMsgValue
	Status StatusMsg
}

type InfoMsg struct {
	Val    string
	Status StatusMsg
}
type UrlDetectionMsgValue struct {
	Urls []string
}
type UrlDetectionMsg struct {
	Val    UrlDetectionMsgValue
	Status StatusMsg
}
type PodMsg struct {
	Val    PodMsgValue
	Status StatusMsg
}

type ContainerMsg struct {
	Val    ContainerMsgValue
	Status StatusMsg
}

type ListUpdateMsg struct {
	Val              []list.Item
	Title            string
	StatusTxt        string
	Status           StatusMsg
	PreSelectedValue string
}

type MainModelChangeFocusMsg struct {
	Val    int
	Status StatusMsg
}
