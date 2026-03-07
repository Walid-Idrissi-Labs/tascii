package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/walid-idrissi-labs/tascii/internal/display"
	"github.com/walid-idrissi-labs/tascii/internal/task"
)

var editCmd = &cobra.Command{
	Use:   "edit [id]",
	Short: "Edit a task's fields",
	Long: `Edit one or more fields of an existing task.
Only the flags you provide will be changed.

Examples:
  tascii edit 3 --title "Updated title"
  tascii edit 3 --priority high --due 2025-01-15
  tascii edit 3 --tag work --tag urgent
  tascii edit 3 --clear-due`,
	Args: cobra.ExactArgs(1),
	RunE: runEdit,
}

var (
	editTitle    string
	editPriority string
	editDue      string
	editNote     string
	editTags     []string
	editClearDue bool
)

func init() {
	editCmd.Flags().StringVar(&editTitle, "title", "", "New title")
	editCmd.Flags().StringVarP(&editPriority, "priority", "p", "", "New priority: 0 (low), 1 (med), or 2 (high); can also use low/med/high")
	editCmd.Flags().StringVarP(&editDue, "due", "d", "", "New due date (YYYY-MM-DD)")
	editCmd.Flags().StringVarP(&editNote, "note", "n", "", "New note (replaces existing)")
	editCmd.Flags().StringArrayVarP(&editTags, "tag", "t", nil, "New tags — replaces all existing tags")
	editCmd.Flags().BoolVar(&editClearDue, "clear-due", false, "Remove the due date entirely")
	rootCmd.AddCommand(editCmd)
}

func runEdit(cmd *cobra.Command, args []string) error {
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID %q — must be a number", args[0])
	}


	titleChanged    := cmd.Flags().Changed("title")
	priorityChanged := cmd.Flags().Changed("priority")
	dueChanged      := cmd.Flags().Changed("due")
	noteChanged     := cmd.Flags().Changed("note")
	tagsChanged     := cmd.Flags().Changed("tag")

	if !titleChanged && !priorityChanged && !dueChanged && !noteChanged && !tagsChanged && !editClearDue {
		return fmt.Errorf("nothing to update — provide at least one flag (e.g. --title, --priority)")
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
		if tasks[i].ID != id {
			continue
		}

		if titleChanged {
			tasks[i].Title = strings.TrimSpace(editTitle)
		}
		if priorityChanged {
			p, err := task.ParsePriority(editPriority)
			if err != nil {
				return err
			}
			tasks[i].Priority = p
		}
		if dueChanged {
			if _, err := time.Parse("2006-01-02", editDue); err != nil {
				return fmt.Errorf("invalid date %q — use YYYY-MM-DD", editDue)
			}
			tasks[i].Due = editDue
		}
		if editClearDue {
			tasks[i].Due = ""
		}
		if noteChanged {
			tasks[i].Note = editNote
		}
		if tagsChanged {
			tasks[i].Tags = editTags
		}

		updated = true
		display.PrintSuccess(fmt.Sprintf("Task #%d updated.", id))
		break
	}

	if !updated {
		return fmt.Errorf("task #%d not found", id)
	}

	return store.Save(tasks)
}
