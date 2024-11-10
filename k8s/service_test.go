package k8s

import (
	"testing"
)

func TestListNamespace(t *testing.T) {
	t.Log("Testinqg the listService Functionality")
}

func TestGetContainer(t *testing.T) {
	k8s := GetService()
	liste, _ := k8s.GetContainerFromPods("things-api-6b678875f6-svsvq", "things")
	t.Log(liste[0].Name)
}
