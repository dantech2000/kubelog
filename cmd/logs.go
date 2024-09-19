package cmd

import (
	"fmt"
	"os"

	lib "github.com/dantech2000/kubelog/lib"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs [container_id]",
	Short: "Display logs for a specific container",
	Long: `Display logs for a specific container. You can filter logs by level using the --level flag.
Supported levels are DEBUG, INFO, WARN, and ERROR.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		podName := args[0]
		namespace, _ := cmd.Flags().GetString("namespace")
		container, _ := cmd.Flags().GetString("container")
		follow, _ := cmd.Flags().GetBool("follow")
		levelStr, _ := cmd.Flags().GetString("level")

		clientset, err := lib.GetKubernetesClient()
		if err != nil {
			fmt.Println("error creating Kubernetes client:", err)
			os.Exit(1)
		}

		logFetcher := lib.NewLogFetcher(clientset, namespace, podName, container, follow)

		filterLevel, err := lib.ParseLogLevel(levelStr)
		if err != nil {
			fmt.Printf("Invalid log level: %v\n", err)
			os.Exit(1)
		}

		logReader, err := logFetcher.GetLogReader()
		if err != nil {
			fmt.Printf("Error fetching logs: %v\n", err)
			os.Exit(1)
		}
		defer logReader.Close()

		if err := lib.FilterAndFormatLogs(logReader, os.Stdout, filterLevel); err != nil {
			fmt.Printf("Error processing logs: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd) // Make sure this line exists and is called only once
	logsCmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace of the pod")
	logsCmd.Flags().StringP("container", "c", "", "Specific container name within the pod (optional if pod has only one container)")
	logsCmd.Flags().BoolP("follow", "f", false, "Follow the log output in real-time (similar to 'tail -f')")
	logsCmd.Flags().StringP("level", "l", "DEBUG", "Filter logs by level (DEBUG, INFO, WARN, ERROR)")
}
