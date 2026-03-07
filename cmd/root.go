package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "tascii",
	Short: "A terminal task manager",
	Long: `tascii — a fast, minimal task manager built for the terminal.

Manage tasks with priorities, deadlines, tags, and notes.
Data is stored locally at: ~/.local/share/tascii/tasks.json`,

	SilenceUsage: true,
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = "1.1.1"
}
