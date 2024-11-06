package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/louislouislouislouis/repr8ducer/ui"
)

var DEFAULT_NAMESPACE = "things"

// func main() {
// 	k8sService := k8s.GetService()
// 	//k8sService.ListNamespace()
// 	//k8sService.ListPodsInNamespace(DEFAULT_NAMESPACE)
// 	err := k8sService.Exec()
// 	fmt.Println(err)
// }

func main() {
	//k8sService := k8s.GetService()
	//itemss, _ := k8sService.ListNamespace()
	// 	items := []list.Item{}
	//
	// 	for _, o := range itemss.Items {
	// 		newItem := item{
	// 			title: o.Name,             // Utilisez le champ que vous voulez pour le tqitre
	// 			desc:  string(o.GetUID()), // Description personnalisée
	// 		}
	// 		items = append(items, newItem)
	// 	}

	listeNameSpace := ui.CreateFakeList(99, "Namespace")
	m := ui.NewModel(listeNameSpace)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
