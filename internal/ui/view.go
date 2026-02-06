package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/flyme2mars/jotcli/internal/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true).Underline(true)
	previewStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).BorderForeground(lipgloss.Color("240"))
)

type editFinishedMsg struct{ err error }

type model struct {
	notes       []database.Note
	cursor      int
	err         error
	quitting    bool
	editingFile string
	editingID   int
}

func InitialModel() model {
	notes, err := database.GetNotes("")
	return model{
		notes:  notes,
		cursor: 0,
		err:    err,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case editFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}

		// Read the updated content
		updatedContent, err := os.ReadFile(m.editingFile)
		if err != nil {
			m.err = err
			return m, nil
		}

		// Save back to database
		err = database.UpdateNote(m.editingID, string(updatedContent))
		if err != nil {
			m.err = err
			return m, nil
		}

		// Clean up
		os.Remove(m.editingFile)
		m.editingFile = ""
		m.editingID = 0

		// Refresh notes
		m.notes, m.err = database.GetNotes("")
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.notes)-1 {
				m.cursor++
			}
		case "e":
			if len(m.notes) > 0 {
				note := m.notes[m.cursor]
				
				// Create temp file
				tmpFile, err := os.CreateTemp("", "jot-*.md")
				if err != nil {
					m.err = err
					return m, nil
				}
				tmpFile.WriteString(note.Content)
				tmpFile.Close()

				m.editingFile = tmpFile.Name()
				m.editingID = note.ID

				// Determine editor
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "vim"
				}

				c := exec.Command(editor, m.editingFile)
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return editFinishedMsg{err}
				})
			}
		case "delete", "x", "backspace":
			if len(m.notes) > 0 {
				note := m.notes[m.cursor]
				err := database.DeleteNote(note.ID)
				if err != nil {
					m.err = err
					return m, nil
				}
				// Refresh notes after deletion
				m.notes, m.err = database.GetNotes("")
				if m.cursor >= len(m.notes) && m.cursor > 0 {
					m.cursor = len(m.notes) - 1
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	if m.quitting {
		return "Bye!\n"
	}

	if len(m.notes) == 0 {
		return "No notes yet. Add one with 'jot add'!\n\n(press q to quit)"
	}

	s := titleStyle.Render("--- Your Notes ---") + "\n\n"

	for i, note := range m.notes {
		cursor := "  "
		// Replace actual newlines and literal \n with spaces for the list view
		displayContent := strings.ReplaceAll(note.Content, "\n", " ")
		displayContent = strings.ReplaceAll(displayContent, "\\n", " ")
		
		line := fmt.Sprintf("[%s] %s", note.Tag, displayContent)
		
		if len(line) > 50 {
			line = line[:47] + "..."
		}

		if m.cursor == i {
			cursor = "> "
			s += selectedStyle.Render(fmt.Sprintf("%s%s", cursor, line)) + "\n"
		} else {
			s += normalStyle.Render(fmt.Sprintf("%s%s", cursor, line)) + "\n"
		}
	}

	// Preview section for the selected note
	if len(m.notes) > 0 {
		selectedNote := m.notes[m.cursor]
		// Ensure literal \n are treated as real newlines for Glamour
		previewContent := strings.ReplaceAll(selectedNote.Content, "\\n", "\n")
		rendered, _ := glamour.Render(previewContent, "dark")
		s += "\n" + previewStyle.Render(rendered)
	}

	s += "\n(up/down: navigate • e: edit • x: delete • q: quit)\n"

	return s
}
