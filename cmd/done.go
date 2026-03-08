package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/walid-idrissi-labs/tascii/internal/display"
	"github.com/walid-idrissi-labs/tascii/internal/task"
)

var doneCmd = &cobra.Command{
	Use:   "done [id... | all]",
	Short: "Mark one, multiple or all tasks as done",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runDone,
}

func init() {
	rootCmd.AddCommand(doneCmd)
}

func runDone(cmd *cobra.Command, args []string) error {

	store, err := task.NewStore()
	if err != nil {
		return err
	}

	tasks, err := store.Load()
	if err != nil {
		return err
	}

	// Handle "tascii done all"
	if len(args) == 1 && args[0] == "all" {
		marked := 0
		for i := range tasks {
			if tasks[i].Status != task.StatusDone {
				tasks[i].Status = task.StatusDone
				display.PrintSuccess(fmt.Sprintf("Task #%d marked as done: %s", tasks[i].ID, tasks[i].Title))
				marked++
			}
		}
		if marked == 0 {
			return fmt.Errorf("no pending tasks to mark as done")
		}
		return store.Save(tasks)
	}

	// Parse numeric IDs
	ids := make([]int, 0, len(args))
	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("invalid task ID %q — must be a number", arg)
		}
		ids = append(ids, id)
	}

	for _, id := range ids {
		found := false
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Status = task.StatusDone
				found = true
				display.PrintSuccess(fmt.Sprintf("Task #%d marked as done: %s", id, tasks[i].Title))
				break
			}
		}
		if !found {
			return fmt.Errorf("task #%d not found", id)
		}
	}

	return store.Save(tasks)
}
