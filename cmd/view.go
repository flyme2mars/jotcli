package cmd

import (
	"fmt"
	"os"

	"github.com/flyme2mars/jot/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Interactive view of your notes",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(ui.InitialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
