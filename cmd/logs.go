package cmd

import (
	"fmt"
	"os"

	lib "github.com/dantech2000/kubelog/lib"
	"github.com/spf13/cobra"
)

// logOptions holds the command options for the logs command
type logOptions struct {
	namespace string
	container string
	follow    bool
	level     string
	podName   string
}

var logsCmd = &cobra.Command{
	Use:   "logs [container_id]",
	Short: "Display logs for a specific container",
	Long: `Display logs for a specific container. You can filter logs by level using the --level flag.
Supported levels are DEBUG, INFO, WARN, and ERROR.`,
	Args: cobra.ExactArgs(1),
	Run:  runLogs,
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().StringP("namespace", "n", "default", "Kubernetes namespace of the pod")
	logsCmd.Flags().StringP("container", "c", "", "Specific container name within the pod (optional if pod has only one container)")
	logsCmd.Flags().BoolP("follow", "f", false, "Follow the log output in real-time (similar to 'tail -f')")
	logsCmd.Flags().StringP("level", "l", "DEBUG", "Filter logs by level (DEBUG, INFO, WARN, ERROR)")
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

	return &logOptions{
		namespace: namespace,
		container: container,
		follow:    follow,
		level:     level,
		podName:   args[0],
	}, nil
}

func runLogs(cmd *cobra.Command, args []string) {
	opts, err := getLogOptions(cmd, args)
	if err != nil {
		fmt.Printf("Error getting command options: %v\n", err)
		os.Exit(1)
	}

	clientset, err := lib.GetKubernetesClient()
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}

	logFetcher := lib.NewLogFetcher(clientset, opts.namespace, opts.podName, opts.container, opts.follow)

	filterLevel, err := lib.ParseLogLevel(opts.level)
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
}
