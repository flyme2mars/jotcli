package cmd

import (
	"fmt"
	"strings"

	"github.com/flyme2mars/jotcli/internal/database"
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
		note := strings.Join(args, " ")
		// Convert literal \n to actual newlines
		note = strings.ReplaceAll(note, "\\n", "\n")
		
		err := database.AddNote(note, tag, priority)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("✅ Note saved: %s\n", note)
	},
}

func init() {
	addCmd.Flags().StringVarP(&tag, "tag", "t", "", "Tag for the note")
	addCmd.Flags().StringVarP(&priority, "priority", "p", "low", "Priority level (low, medium, high)")
	rootCmd.AddCommand(addCmd)
}
