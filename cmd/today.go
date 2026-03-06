package cmd


import (
	"github.com/spf13/cobra"
	"github.com/walid-idrissi-labs/tascii/internal/display"
	"github.com/walid-idrissi-labs/tascii/internal/task"
)

var todayCmd = &cobra.Command{
	Use:   "today",
	Short: "Show tasks due today and overdue",
	Args:  cobra.NoArgs,
	RunE:  runToday,
}

func init() {
	rootCmd.AddCommand(todayCmd)
}

func runToday(cmd *cobra.Command, args []string) error {
	store, err := task.NewStore()
	if err != nil {
		return err
	}

	tasks, err := store.Load()
	if err != nil {
		return err
	}


	urgent := make([]task.Task, 0)
	for _, t := range tasks {
		if t.Status != task.StatusDone && (t.IsDueToday() || t.IsOverdue()) {
			urgent = append(urgent, t)
		}
	}


	urgent = task.Sort(urgent, "due")

	display.PrintToday(urgent)
	return nil
}
