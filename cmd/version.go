package cmd

import (
	"fmt"

	lib "github.com/dantech2000/kubelog/lib"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of kubelog",
	Long:  `All software has versions. This is kubelog's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(lib.CurrentVersion.FullString())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
