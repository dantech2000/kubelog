package cmd

import (
	"fmt"
	"os"

	"github.com/dantech2000/kubelog/pkg/kubernetes"
	"github.com/spf13/cobra"
)

// logOptions holds the command options for the logs command
type logOptions struct {
	namespace string
	container string
	follow    bool
	level     string
	podName   string
	previous  bool
}

var logsCmd = &cobra.Command{
	Use:   "logs [container_id]",
	Short: "Display logs for a specific container",
	Long: `Display logs for a specific container. You can filter logs by level using the --level flag.
Supported levels are DEBUG, INFO, WARN, and ERROR.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runLogs(cmd, args); err != nil {
			fmt.Printf("Error running logs command: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().StringP("namespace", "n", "", "Kubernetes namespace (defaults to current context's namespace)")
	logsCmd.Flags().StringP("container", "c", "", "Specific container name within the pod (optional if pod has only one container)")
	logsCmd.Flags().BoolP("follow", "f", false, "Follow the log output in real-time (similar to 'tail -f')")
	logsCmd.Flags().StringP("level", "l", "DEBUG", "Filter logs by level (DEBUG, INFO, WARN, ERROR)")
	logsCmd.Flags().BoolP("previous", "p", false, "Get previous terminated container logs")
}

func getLogOptions(cmd *cobra.Command, args []string) (*logOptions, error) {
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		return nil, fmt.Errorf("error getting namespace flag: %v", err)
	}

	container, err := cmd.Flags().GetString("container")
	if err != nil {
		return nil, fmt.Errorf("error getting container flag: %v", err)
	}

	follow, err := cmd.Flags().GetBool("follow")
	if err != nil {
		return nil, fmt.Errorf("error getting follow flag: %v", err)
	}

	level, err := cmd.Flags().GetString("level")
	if err != nil {
		return nil, fmt.Errorf("error getting level flag: %v", err)
	}

	previous, err := cmd.Flags().GetBool("previous")
	if err != nil {
		return nil, fmt.Errorf("error getting previous flag: %v", err)
	}

	return &logOptions{
		namespace: namespace,
		container: container,
		follow:    follow,
		level:     level,
		podName:   args[0],
		previous:  previous,
	}, nil
}

func runLogs(cmd *cobra.Command, args []string) error {
	options, err := getLogOptions(cmd, args)
	if err != nil {
		return err
	}

	clientset, contextNamespace, err := kubernetes.GetKubernetesClient()
	if err != nil {
		return fmt.Errorf("error getting kubernetes client: %v", err)
	}

	// Use context namespace if no namespace is specified
	if options.namespace == "" {
		options.namespace = contextNamespace
	}

	// Create log fetcher with the new interface
	logFetcher := kubernetes.NewLogFetcher(
		clientset,
		options.namespace,
		options.podName,
		options.follow,
		options.previous,
		os.Stdout,
	)

	// Get logs using the new method
	err = logFetcher.GetLogs()
	if err != nil {
		return fmt.Errorf("error fetching logs: %v", err)
	}

	return nil
}
