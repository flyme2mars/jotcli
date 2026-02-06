package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/flyme2mars/jotcli/internal/database"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [id]",
	Short: "Edit a note in your default editor",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: Invalid note ID")
			return
		}

		note, err := database.GetNoteByID(id)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if note == nil {
			fmt.Printf("Error: Note with ID %d not found\n", id)
			return
		}

		// Create a temporary file
		tmpFile, err := os.CreateTemp("", "jot-*.md")
		if err != nil {
			fmt.Printf("Error: Could not create temp file: %v\n", err)
			return
		}
		defer os.Remove(tmpFile.Name())

		// Write the current content to the temp file
		if _, err := tmpFile.WriteString(note.Content); err != nil {
			fmt.Printf("Error: Could not write to temp file: %v\n", err)
			return
		}
		tmpFile.Close()

		// Determine which editor to use
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim" // Fallback to vim
		}

		// Open the editor
		editProcess := exec.Command(editor, tmpFile.Name())
		editProcess.Stdin = os.Stdin
		editProcess.Stdout = os.Stdout
		editProcess.Stderr = os.Stderr

		if err := editProcess.Run(); err != nil {
			fmt.Printf("Error: Editor failed: %v\n", err)
			return
		}

		// Read the updated content
		updatedContent, err := os.ReadFile(tmpFile.Name())
		if err != nil {
			fmt.Printf("Error: Could not read updated file: %v\n", err)
			return
		}

		// Save back to database
		err = database.UpdateNote(id, string(updatedContent))
		if err != nil {
			fmt.Printf("Error saving note: %v\n", err)
			return
		}

		fmt.Printf("✅ Note %d updated!\n", id)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}