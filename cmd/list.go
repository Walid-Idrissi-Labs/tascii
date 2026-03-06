package cmd


import (
	"github.com/spf13/cobra"
	"github.com/walid-idrissi-labs/tascii/internal/display"
	"github.com/walid-idrissi-labs/tascii/internal/task"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long: `List all tasks in a table. Use flags to filter and sort.

Examples:
  tascii list
  tascii list --tag work
  tascii list --sort priority
  tascii list --filter in-progress
  tascii list -t work -s due`,


	Aliases: []string{"ls"},


	Args: cobra.NoArgs,
	RunE: runList,
}

var (
	listTag    string 
	listSort   string 
	listFilter string 
)

func init() {
	listCmd.Flags().StringVarP(&listTag, "tag", "t", "", "Filter by tag")
	listCmd.Flags().StringVarP(&listSort, "sort", "s", "id", "Sort by: id, priority, due, created")
	listCmd.Flags().StringVarP(&listFilter, "filter", "f", "", "Filter by status: todo, in-progress, done")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	store, err := task.NewStore()
	if err != nil {
		return err
	}

	tasks, err := store.Load()
	if err != nil {
		return err
	}

	tasks = task.Filter(tasks, listTag, listFilter)
	tasks = task.Sort(tasks, listSort)

	display.PrintTable(tasks)
	return nil
}
