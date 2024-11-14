package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Descripable interface {
	GetName() string
	GetUID() string
}

func createDisplayedListFromMetadata[T any](slice []T, getMetadata func(T) Descripable) []list.Item {
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

func CreateFakeList(numItems int, descr string) []list.Item {
	items := make([]list.Item, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = DisplayedItem{
			title: fmt.Sprintf("%s %d", descr, i+1),
			desc:  fmt.Sprintf("Description for %s %d", descr, i+1),
		}
	}
	return items
}

type descipableMeta struct {
	meta metaV1.Object
}

type descipableContainer struct {
	v1.Container
}

func (container descipableMeta) GetUID() string {
	return string(container.meta.GetUID())
}

func (container descipableMeta) GetName() string {
	return string(container.meta.GetName())
}

func (container descipableContainer) GetName() string {
	return container.Name
}

func (container descipableContainer) GetUID() string {
	return container.Image
}
