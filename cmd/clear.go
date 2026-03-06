package cmd
import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/walid-idrissi-labs/tascii/internal/display"
	"github.com/walid-idrissi-labs/tascii/internal/task"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all completed tasks",
	Args:  cobra.NoArgs,
	RunE:  runClear,
}

var clearForce bool 

func init() {
	clearCmd.Flags().BoolVarP(&clearForce, "force", "f", false, "Skip confirmation prompt")
	rootCmd.AddCommand(clearCmd)
}

func runClear(cmd *cobra.Command, args []string) error {
	store, err := task.NewStore()
	if err != nil {
		return err
	}

	tasks, err := store.Load()
	if err != nil {
		return err
	}


	count := 0
	for _, t := range tasks {
		if t.Status == task.StatusDone {
			count++
		}
	}

	if count == 0 {
		display.PrintInfo("No completed tasks to clear.")
		return nil
	}

	if !clearForce {
		fmt.Printf("\n  Delete %d completed task(s)? [y/N] ", count)
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


	remaining := make([]task.Task, 0)
	for _, t := range tasks {
		if t.Status != task.StatusDone {
			remaining = append(remaining, t)
		}
	}

	if err := store.Save(remaining); err != nil {
		return err
	}

	display.PrintSuccess(fmt.Sprintf("Cleared %d completed task(s).", count))
	return nil
}
