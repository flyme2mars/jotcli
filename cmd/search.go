package cmd

import (
	"fmt"
	"strings"

	"github.com/flyme2mars/jotcli/internal/database"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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
			cmd.Printf("Error searching notes: %v\n", err)
			return
		}
		defer rows.Close()

		var notes []database.Note
		for rows.Next() {
			var n database.Note
			err := rows.Scan(&n.ID, &n.Content, &n.Tag, &n.Priority, &n.CreatedAt)
			if err != nil {
				cmd.Printf("Error scanning result: %v\n", err)
				return
			}
			notes = append(notes, n)
		}

		if len(notes) == 0 {
			cmd.Printf("No notes found matching '%s'\n", query)
			return
		}

		headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true).Padding(0, 1)
		cellStyle := lipgloss.NewStyle().Padding(0, 1)
		borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

		rowsTable := [][]string{}
		for _, n := range notes {
			content := strings.ReplaceAll(n.Content, "\n", " ")
			rowsTable = append(rowsTable, []string{
				fmt.Sprintf("%d", n.ID),
				content,
				n.Tag,
				n.Priority,
				n.CreatedAt.Format("2006-01-02"),
			})
		}

		t := table.New().
			Border(lipgloss.NormalBorder()).
			BorderStyle(borderStyle).
			StyleFunc(func(row, col int) lipgloss.Style {
				if row == table.HeaderRow {
					return headerStyle
				}
				return cellStyle
			}).
			Headers("ID", "Note", "Tag", "Priority", "Created").
			Rows(rowsTable...)

		cmd.Println(t.Render())
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}