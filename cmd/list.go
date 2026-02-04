package cmd

import (
	"fmt"
	"os"

	"github.com/flyme2mars/jotcli/internal/database"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listTag string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Run: func(cmd *cobra.Command, args []string) {
		notes, err := database.GetNotes(listTag)
		if err != nil {
			fmt.Printf("Error retrieving notes: %v\n", err)
			return
		}

		if len(notes) == 0 {
			fmt.Println("No notes found.")
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header("ID", "Note", "Tag", "Priority", "Created At")

		for _, n := range notes {
			table.Append([]string{
				fmt.Sprintf("%d", n.ID),
				n.Content,
				n.Tag,
				n.Priority,
				n.CreatedAt.Format("2006-01-02 15:04"),
			})
		}
		table.Render()
	},
}

func init() {
	listCmd.Flags().StringVarP(&listTag, "tag", "t", "", "Filter notes by tag")
	rootCmd.AddCommand(listCmd)
}
