package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/walid-idrissi-labs/tascii/internal/display"
	"github.com/walid-idrissi-labs/tascii/internal/task"
)

var viewCmd = &cobra.Command{
	Use:   "view [id]",
	Short: "View full details of a task",
	Args:  cobra.ExactArgs(1),
	RunE:  runView,
}

func init() {
	rootCmd.AddCommand(viewCmd)
}

func runView(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID %q — must be a number", args[0])
	}

	store, err := task.NewStore()
	if err != nil {
		return err
	}

	tasks, err := store.Load()
	if err != nil {
		return err
	}

	for _, t := range tasks {
		if t.ID == id {
			display.PrintDetail(t)
			return nil
		}
	}

	return fmt.Errorf("task #%d not found", id)
}
