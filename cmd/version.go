package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const AppVersion = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of TaskManagerCLI",
	Long:  `All software has versions. This is TaskManagerCLI's!`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("TaskManagerCLI Version: %s\n", AppVersion)
	},
}

func init() {
	AddCommand(versionCmd)
}
