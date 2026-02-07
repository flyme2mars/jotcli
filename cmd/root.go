package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/flyme2mars/jotcli/internal/database"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jotcli",
	Short: "jotcli is a CLI tool for capturing thoughts",
	Long:  `A quick and efficient way to capture notes, tag them, and view them in your terminal.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// This runs BEFORE any subcommand
		err := database.InitDB()
		if err != nil {
			fmt.Printf("Error initializing database: %v\n", err)
			os.Exit(1)
		}
	},
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

// SetOutput allows redirecting the CLI output (useful for tests)
func SetOutput(w io.Writer) {
	rootCmd.SetOut(w)
	rootCmd.SetErr(w)
}
