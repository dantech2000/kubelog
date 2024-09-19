package cmd

import (
	"fmt"
	"os"

	lib "github.com/dantech2000/kubelog/lib"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var containersCmd = &cobra.Command{
	Use:   "containers [pod-name]",
	Short: "List containers in a Kubernetes pod",
	Long: `List all containers within a specified Kubernetes pod.
This command provides a formatted output of container names for the given pod,
including the total count of containers.

Example usage:
  kubelog containers my-pod -n my-namespace
  kubelog containers my-pod -n my-namespace --output json
  kubelog containers my-pod -n my-namespace -o yaml
  kubelog containers my-pod -n my-namespace -o posix`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podName := args[0]
		namespace, _ := cmd.Flags().GetString("namespace")
		outputFormat, _ := cmd.Flags().GetString("output")

		clientset, err := lib.GetKubernetesClient()
		if err != nil {
			color.Red("Error creating Kubernetes client: %v", err)
			os.Exit(1)
		}

		containers, err := lib.ListContainers(clientset, namespace, podName)
		if err != nil {
			color.Red("Error listing containers: %v", err)
			os.Exit(1)
		}

		formatter := lib.NewOutputFormatter(podName, namespace, containers)
		output, err := formatter.FormatOutput(outputFormat)
		if err != nil {
			color.Red("Error formatting output: %v", err)
			os.Exit(1)
		}

		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(containersCmd)
	containersCmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace of the pod")
	containersCmd.Flags().StringP("output", "o", "", "Output format: json, yaml, or posix")
}
