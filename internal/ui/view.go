package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/flyme2mars/jotcli/internal/database"
	"github.com/charmbracelet/bubbles/textinput"
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
	promptStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
)

type editFinishedMsg struct{ err error }

type mode int

const (
	modeList mode = iota
	modeInput
)

type model struct {
	notes       []database.Note
	cursor      int
	err         error
	quitting    bool
	editingFile string
	editingID   int
	
	// New fields for input mode
	mode      mode
	textInput textinput.Model
}

func InitialModel() model {
	notes, err := database.GetNotes("")
	
	ti := textinput.New()
	ti.Placeholder = "Enter your thought..."
	ti.CharLimit = 250
	ti.Width = 50

	return model{
		notes:     notes,
		cursor:    0,
		err:       err,
		mode:      modeList,
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle Input Mode
	if m.mode == modeInput {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				if m.textInput.Value() != "" {
					err := database.AddNote(m.textInput.Value(), "inbox", "low")
					if err != nil {
						m.err = err
						return m, nil
					}
					m.notes, m.err = database.GetNotes("")
				}
				m.mode = modeList
				m.textInput.Reset()
				return m, nil
			case "esc":
				m.mode = modeList
				m.textInput.Reset()
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	// Handle List Mode
	switch msg := msg.(type) {
	case editFinishedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		updatedContent, err := os.ReadFile(m.editingFile)
		if err != nil {
			m.err = err
			return m, nil
		}
		err = database.UpdateNote(m.editingID, string(updatedContent))
		if err != nil {
			m.err = err
			return m, nil
		}
		os.Remove(m.editingFile)
		m.editingFile = ""
		m.editingID = 0
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
		case "n":
			m.mode = modeInput
			m.textInput.Focus()
			return m, textinput.Blink
		case "e":
			if len(m.notes) > 0 {
				note := m.notes[m.cursor]
				tmpFile, err := os.CreateTemp("", "jot-*.md")
				if err != nil {
					m.err = err
					return m, nil
				}
				tmpFile.WriteString(note.Content)
				tmpFile.Close()
				m.editingFile = tmpFile.Name()
				m.editingID = note.ID
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

	// View for Input Mode
	if m.mode == modeInput {
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s",
			titleStyle.Render("--- New Note ---"),
			m.textInput.View(),
			"(esc to cancel • enter to save)",
		)
	}

	// View for List Mode
	if len(m.notes) == 0 {
		return "No notes yet. Press 'n' to add one!\n\n(press q to quit)"
	}

	s := titleStyle.Render("--- Your Notes ---") + "\n\n"

	for i, note := range m.notes {
		cursor := "  "
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

	if len(m.notes) > 0 {
		selectedNote := m.notes[m.cursor]
		previewContent := strings.ReplaceAll(selectedNote.Content, "\\n", "\n")
		rendered, _ := glamour.Render(previewContent, "dark")
		s += "\n" + previewStyle.Render(rendered)
	}

	s += "\n(n: new • e: edit • x: delete • q: quit)\n"
	return s
}