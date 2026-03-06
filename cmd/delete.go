package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/walid-idrissi-labs/tascii/internal/display"
	"github.com/walid-idrissi-labs/tascii/internal/task"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [id]",
	Short:   "Delete a task",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE:    runDelete,
}

var deleteForce bool 

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Skip confirmation prompt")
	rootCmd.AddCommand(deleteCmd)
}

func runDelete(cmd *cobra.Command, args []string) error {
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


	var found *task.Task
	for i := range tasks {
		if tasks[i].ID == id {
			found = &tasks[i]
			break
		}
	}

	if found == nil {
		return fmt.Errorf("task #%d not found", id)
	}


	if !deleteForce {
		fmt.Printf("\n  Delete task #%d: %q? [y/N] ", id, found.Title)


		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))


		if input != "y" && input != "yes" {
			fmt.Println()
			display.PrintInfo("Cancelled.")
			return nil
		}
		fmt.Println()
	}


	remaining := make([]task.Task, 0, len(tasks)-1)
	for _, t := range tasks {
		if t.ID != id {
			remaining = append(remaining, t)
		}
	}

	if err := store.Save(remaining); err != nil {
		return err
	}

	display.PrintSuccess(fmt.Sprintf("Task #%d deleted.", id))
	return nil
}
