package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubelog",
	Short: "Kubelog - A CLI tool for enhanced Kubernetes log viewing",
	Long: `Kubelog is a command-line interface tool designed to simplify and enhance
the process of viewing Kubernetes pod logs.

It provides features such as:
- Fetching and displaying logs from specific pods and containers
- Real-time log following
- Listing containers within a pod
- Color-coded output for improved readability

Use "kubelog [command] --help" for more information about a command.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you can define flags and configuration settings that are global to all commands.
	// For example, setting a default namespace.
	rootCmd.PersistentFlags().StringP("namespace", "n", "default", "Kubernetes namespace")
	rootCmd.AddCommand(versionCmd)
}
