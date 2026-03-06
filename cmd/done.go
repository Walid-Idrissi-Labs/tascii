package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/walid-idrissi-labs/tascii/internal/display"
	"github.com/walid-idrissi-labs/tascii/internal/task"
)

var doneCmd = &cobra.Command{
	Use:   "done [id]",
	Short: "Mark a task as done",
	Args:  cobra.ExactArgs(1), 
	RunE:  runDone,
}

func init() {
	rootCmd.AddCommand(doneCmd)
}

func runDone(cmd *cobra.Command, args []string) error {

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


	updated := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Status = task.StatusDone
			updated = true
			display.PrintSuccess(fmt.Sprintf("Task #%d marked as done: %s", id, tasks[i].Title))
			break
		}
	}

	if !updated {
		return fmt.Errorf("task #%d not found", id)
	}

	return store.Save(tasks)
}
