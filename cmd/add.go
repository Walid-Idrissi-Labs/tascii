package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/walid-idrissi-labs/tascii/internal/display"
	"github.com/walid-idrissi-labs/tascii/internal/task"
)


var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new task",
	Long: `Add a new task with optional priority, due date, tags, and a note.

Examples:
  tascii add "Fix login bug"
  tascii add "Fix login bug" --priority high --due 2024-12-25
  tascii add "Fix login bug" -p 2 -d 2024-12-25 -t work -n "Affects OAuth flow"
  tascii add --tag "urgent"`,


	Args: cobra.ArbitraryArgs,


	RunE: runAdd,
}


var (
	addPriority string   
	addDue      string   
	addTags     []string 
	addNote     string   
)

func init() {
	addCmd.Flags().StringVarP(&addPriority, "priority", "p", "" , "Priority: 0 (low), 1 (med), or 2 (high); can also use low/med/high")
	addCmd.Flags().StringVarP(&addDue, "due", "d", "", "Due date in YYYY-MM-DD format")
	addCmd.Flags().StringArrayVarP(&addTags, "tag", "t", []string{}, "Tag (can repeat: -t work -t urgent)")
	addCmd.Flags().StringVarP(&addNote, "note", "n", "", "Optional note or description")

	rootCmd.AddCommand(addCmd)
}


func runAdd(cmd *cobra.Command, args []string) error {
	title := strings.TrimSpace(strings.Join(args, " "))
	if title == "" {
		title = "New Reminder"
	}

	var priority int
	if addPriority != "" {
		var err error
		priority, err = task.ParsePriority(addPriority)
		if err != nil {
			return err
		}
	}
	

	if addDue != "" {
		if _, err := time.Parse("2006-01-02", addDue); err != nil {
			return fmt.Errorf("invalid due date %q — use YYYY-MM-DD (e.g. 2024-12-25)", addDue)
		}
	}

	store, err := task.NewStore()
	if err != nil {
		return err
	}

	tasks, err := store.Load()
	if err != nil {
		return err
	}


	t := task.Task{
		ID:        store.NextID(tasks),
		Title:     title,
		Priority:  priority,
		Status:    task.StatusTodo,
		Due:       addDue,
		Tags:      addTags,
		Note:      addNote,
		CreatedAt: time.Now(),
	}


	tasks = append(tasks, t)




	if err := store.Save(tasks); err != nil {
		return err
	}

	display.PrintSuccess(fmt.Sprintf("Task #%d added: %s", t.ID, t.Title))
	return nil

}
