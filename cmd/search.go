package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/flyme2mars/jotcli/internal/database"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for notes containing a query string",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]
		
		sqlQuery := `SELECT id, content, tag, priority, created_at FROM notes WHERE content LIKE ? ORDER BY created_at DESC`
		rows, err := database.DB.Query(sqlQuery, "%"+query+"%")
		if err != nil {
			fmt.Printf("Error searching notes: %v\n", err)
			return
		}
		defer rows.Close()

		var notes []database.Note
		for rows.Next() {
			var n database.Note
			err := rows.Scan(&n.ID, &n.Content, &n.Tag, &n.Priority, &n.CreatedAt)
			if err != nil {
				fmt.Printf("Error scanning result: %v\n", err)
				return
			}
			notes = append(notes, n)
		}

		if len(notes) == 0 {
			fmt.Printf("No notes found matching '%s'\n", query)
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
	rootCmd.AddCommand(searchCmd)
}