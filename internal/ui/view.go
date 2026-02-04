package ui

import (
	"fmt"

	"github.com/flyme2mars/jot/internal/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
)

type model struct {
	notes    []database.Note
	cursor   int
	err      error
	quitting bool
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

	s := "--- Your Notes ---\n\n"

	for i, note := range m.notes {
		cursor := "  "
		line := fmt.Sprintf("[%s] %s", note.Tag, note.Content)
		
		if m.cursor == i {
			cursor = "> "
			s += selectedStyle.Render(fmt.Sprintf("%s%s", cursor, line)) + "\n"
		} else {
			s += normalStyle.Render(fmt.Sprintf("%s%s", cursor, line)) + "\n"
		}
	}

	s += "\n(up/down: navigate • x: delete • q: quit)\n"

	return s
}
