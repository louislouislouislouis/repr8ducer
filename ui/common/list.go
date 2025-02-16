package common

import (
	"github.com/charmbracelet/bubbles/list"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DiplayableItemList interface {
	GetName() string
	GetUID() string
}

type DisplayedItem struct {
	title, desc string
	isSelected  bool
}

func (i DisplayedItem) Title() string       { return i.title }
func (i DisplayedItem) Description() string { return i.desc }
func (i DisplayedItem) FilterValue() string { return i.title }
func (i DisplayedItem) IsSelected() bool    { return i.isSelected }

type DisplayableMeta struct {
	Meta metaV1.Object
}

type DisplayableContainer struct {
	v1.Container
}

func (m DisplayableMeta) GetUID() string {
	return string(m.Meta.GetUID())
}

func (m DisplayableMeta) GetName() string {
	return m.Meta.GetName()
}

func (container DisplayableContainer) GetName() string { return container.Name }

func (container DisplayableContainer) GetUID() string { return container.Image }

func CreateDisplayedListFromMetadata[T any](
	slice []T,
	getMetadata func(T) DiplayableItemList,
) []list.Item {
	items := make([]list.Item, len(slice))
	for i, itemm := range slice {
		meta := getMetadata(itemm)
		newItem := DisplayedItem{
			title: meta.GetName(),
			desc:  meta.GetUID(),
		}
		items[i] = newItem
	}
	return items
}
