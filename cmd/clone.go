package cmd

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/louislouislouislouis/repr8ducer/k8s"
	"github.com/louislouislouislouis/repr8ducer/ui/mainmodel"
	"github.com/louislouislouislouis/repr8ducer/utils"
)

var (
	namespace    string
	podName      string
	container    string
	overrideUrls bool
)

var cloneCmd = &cobra.Command{
	Use:   "reproduce",
	Short: "Reproduce specific pod",
	Long:  "Copy the specific docker command in your keyboard",
	Run: func(cmd *cobra.Command, args []string) {
		// Directly output command if all args are here, and skip url overwrite
		if namespace != "" && podName != "" && container != "" && overrideUrls {
			command, err := k8s.NewDefaultGenerator(k8s.GetService()).PodToContainer(
				namespace,
				podName,
				context.TODO(),
			)
			// TODO, add better error handling, by logging
			if err != nil {
				utils.Log.Error().Msg(err.Error())
			}
			fmt.Println(command.GetCommand())
			return
		}
		// Otherwise return cli
		runCli(namespace, podName, container, overrideUrls)
	},
}

func init() {
	cloneCmd.PersistentFlags().
		StringVarP(&namespace, "namespace", "n", "", "Namespace to work with")
	cloneCmd.PersistentFlags().
		StringVarP(&podName, "podName", "p", "", "Podname to work replicate")
	cloneCmd.PersistentFlags().
		StringVarP(&container, "container", "c", "", "Container to work replicate")
	cloneCmd.Flags().BoolVar(&overrideUrls, "no-url-overwrite", false, "Override Urls in generated file")
	cloneCmd.RegisterFlagCompletionFunc("namespace", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		namespacesList, err := k8s.GetService().ListNamespace(context.Background())
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}
		namespaces := make([]string, len(namespacesList.Items))
		for idx, namespace := range namespacesList.Items {
			namespaces[idx] = namespace.Name
		}

		return namespaces, cobra.ShellCompDirectiveDefault
	})

	cloneCmd.RegisterFlagCompletionFunc("podName", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		flags := cmd.Flags()
		namespace, err := flags.GetString("namespace")
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}
		podList, err := k8s.GetService().ListPodsInNamespace(namespace, context.TODO())
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}

		pods := make([]string, len(podList.Items))
		for idx, pod := range podList.Items {
			pods[idx] = pod.Name
		}
		return pods, cobra.ShellCompDirectiveDefault
	})

	cloneCmd.RegisterFlagCompletionFunc("container", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		flags := cmd.Flags()
		namespace, err := flags.GetString("namespace")
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}
		podName, err := flags.GetString("podName")
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}
		containers, err := k8s.GetService().GetContainerFromPods(namespace, podName, context.TODO())
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}

		containersCompletion := make([]string, len(containers))
		for idx, container := range containers {
			containersCompletion[idx] = container.Name
		}
		return containersCompletion, cobra.ShellCompDirectiveDefault
	})
	rootCmd.AddCommand(cloneCmd)
}

func runCli(namespace, pod, container string, skipUrlOverwrite bool) {
	p := tea.NewProgram(
		mainmodel.NewMainModel(k8s.GetService(), mainmodel.MainModelConfig{
			Pod:              pod,
			Namespace:        namespace,
			Container:        container,
			SkipUrlOverwrite: skipUrlOverwrite,
		}),
		tea.WithAltScreen(),
	)
	_, err := p.Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
