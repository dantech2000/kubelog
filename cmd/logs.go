package cmd

import (
	"fmt"
	"os"

	"github.com/dantech2000/kubelog/utils"

	"github.com/spf13/cobra"
	// ... other imports ...
)

var logsCmd = &cobra.Command{
	Use:   "logs [pod-name]",
	Short: "Fetch and display logs from a Kubernetes pod",
	Long: `Fetch and display logs from a specified Kubernetes pod.
This command allows you to view logs from a pod in real-time or as a one-time fetch.
You can specify the namespace, container, and whether to follow the logs.

Example usage:
  kubelog logs my-pod -n my-namespace -c my-container -f`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podName := args[0]
		namespace, _ := cmd.Flags().GetString("namespace")
		container, _ := cmd.Flags().GetString("container")
		follow, _ := cmd.Flags().GetBool("follow")

		clientset, err := utils.GetKubernetesClient()
		if err != nil {
			fmt.Println("error creating Kubernetes client:", err)
			os.Exit(1)
		}

		logFetcher := utils.NewLogFetcher(clientset, namespace, podName, container, follow)
		if err := logFetcher.FetchLogs(); err != nil {
			fmt.Printf("Error fetching logs: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace of the pod")
	logsCmd.Flags().StringP("container", "c", "", "Specific container name within the pod (optional if pod has only one container)")
	logsCmd.Flags().BoolP("follow", "f", false, "Follow the log output in real-time (similar to 'tail -f')")
}

// ... rest of the file remains the same ...
