package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createDisplayedListFromMetadata[T any](slice []T, getMetadata func(T) metaV1.Object) []list.Item {
	items := make([]list.Item, len(slice))
	for i, itemm := range slice {
		meta := getMetadata(itemm)
		newItem := DisplayedItem{
			title: meta.GetName(),
			desc:  string(meta.GetUID()),
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
