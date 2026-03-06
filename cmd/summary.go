package cmd

import (
	"github.com/spf13/cobra"
	"github.com/walid-idrissi-labs/tascii/internal/display"
	"github.com/walid-idrissi-labs/tascii/internal/task"
)

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show a summary of all tasks",
	Args:  cobra.NoArgs,
	RunE:  runSummary,
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}

func runSummary(cmd *cobra.Command, args []string) error {
	store, err := task.NewStore()
	if err != nil {
		return err
	}

	tasks, err := store.Load()
	if err != nil {
		return err
	}

	display.PrintSummary(tasks)
	return nil
}
