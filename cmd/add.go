package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	tag      string
	priority string
)

var addCmd = &cobra.Command{
	Use:   "add [note]",
	Short: "Add a new note",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		note := args[0]
		fmt.Printf("Note: %s\n", note)
		if tag != "" {
			fmt.Printf("Tag: %s\n", tag)
		}
		if priority != "" {
			fmt.Printf("Priority: %s\n", priority)
		}
		fmt.Println("Note saved (mock)!")
	},
}

func init() {
	addCmd.Flags().StringVarP(&tag, "tag", "t", "", "Tag for the note")
	addCmd.Flags().StringVarP(&priority, "priority", "p", "low", "Priority level (low, medium, high)")
	rootCmd.AddCommand(addCmd)
}
