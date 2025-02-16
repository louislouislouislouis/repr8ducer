package k8s

import (
	"context"
	"testing"
)

func TestListNamespace(t *testing.T) {
	t.Log("Testinqg the listService Functionality")
}

func TestGetContainer(t *testing.T) {
	k8s := GetService()
	liste, _ := k8s.ListNamespace(context.Background())
	t.Log(liste.Items)
}
