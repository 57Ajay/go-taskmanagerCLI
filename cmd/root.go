package cmd

import (
	// "fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "tmgr",
	Short: "A simple CLI task manager and personal assistant",
	Long: `TaskManagerCLI helps you manage tasks, set reminders,
and store quick notes directly from your terminal.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
