package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jot",
	Short: "jot is a CLI tool for capturing thoughts",
	Long:  `A quick and efficient way to capture notes, tag them, and view them in your terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		// This runs when no subcommands are provided
		fmt.Println("Welcome to jot! Use 'jot add' to save a note.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
