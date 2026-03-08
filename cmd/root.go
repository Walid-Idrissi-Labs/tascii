package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/walid-idrissi-labs/tascii/internal/display"
)


var rootCmd = &cobra.Command{
	Use:   "tascii",
	Short: "A terminal task manager",
	Long: `tascii — a fast, minimal task manager built for the terminal.

Manage tasks with priorities, deadlines, tags, and notes.
Data is stored locally at: ~/.local/share/tascii/tasks.json`,

	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		display.PrintMuted("Hi! I'm Tasky, a task manager for the CLI.")
		display.PrintMuted("Made by Walid — check out my profile on GitHub: https://github.com/walid-idrissi-labs")
		display.PrintMuted("")
		display.PrintMuted("To get started, run: tascii -h")
	},
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

var version = "dev" //*default

func init() {
	//* linker to repo 
	rootCmd.Version = version 
}
