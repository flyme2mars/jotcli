package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/flyme2mars/jotcli/internal/database"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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

		// Get terminal width
		width, _, err := term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			width = 80 // Fallback
		}

		// Define styles
		headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true).Padding(0, 1)
		cellStyle := lipgloss.NewStyle().Padding(0, 1)
		borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

		// Prepare data and calculate available space for the "Note" column
		// Fixed widths for ID (4), Tag (12), Priority (10), Created At (18) + Borders
		reservedWidth := 4 + 12 + 10 + 18 + 10 
		noteWidth := width - reservedWidth
		if noteWidth < 20 {
			noteWidth = 20 // Minimum width for the note
		}

		rows := [][]string{}
		for _, n := range notes {
			// Clean up newlines for the table view
			content := strings.ReplaceAll(n.Content, "\n", " ")
			content = strings.ReplaceAll(content, "\\n", " ")

			// Truncate if too long
			if len(content) > noteWidth {
				content = content[:noteWidth-3] + "..."
			}

			rows = append(rows, []string{
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
			Rows(rows...)

		fmt.Println(t.Render())
	},
}

func init() {
	listCmd.Flags().StringVarP(&listTag, "tag", "t", "", "Filter notes by tag")
	rootCmd.AddCommand(listCmd)
}