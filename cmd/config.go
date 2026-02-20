package cmd

import (
	"fmt"

	"github.com/flyme2mars/jotcli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("--- jotcli Configuration ---")
		fmt.Printf("Database Path: %s\n", config.GetDBPath())
		fmt.Printf("Editor:        %s\n", config.GetEditor())
		fmt.Println("\nYou can override these by creating a ~/.jotcli.yaml file")
		fmt.Println("or by setting JOT_DATABASE and EDITOR environment variables.")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
