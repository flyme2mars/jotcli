package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/flyme2mars/jotcli/internal/database"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles
	titleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true).Underline(true)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	previewStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).BorderForeground(lipgloss.Color("240"))
	
	// Textarea styling
	textAreaStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	// Status Bar Style
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"}).
			MarginTop(1)
	
	statusKey = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Padding(0, 1).
			MarginRight(1)
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
	
	mode     mode
	textArea textarea.Model
}

func InitialModel() model {
	notes, err := database.GetNotes("")
	
	ta := textarea.New()
	ta.Placeholder = "What's on your mind?..."
	ta.SetWidth(60)
	ta.SetHeight(10)
	
	// Styling the textarea
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle() // Remove that "hard" gray highlight
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	ta.ShowLineNumbers = false
	
	return model{
		notes:    notes,
		cursor:   0,
		err:      err,
		mode:     modeList,
		textArea: ta,
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
			case "ctrl+s": // Save note
				content := strings.TrimSpace(m.textArea.Value())
				if content != "" {
					err := database.AddNote(content, "inbox", "low")
					if err != nil {
						m.err = err
						return m, nil
					}
					m.notes, m.err = database.GetNotes("")
				}
				m.mode = modeList
				m.textArea.Reset()
				return m, nil
			case "esc":
				m.mode = modeList
				m.textArea.Reset()
				return m, nil
			}
		}

		var cmd tea.Cmd
		m.textArea, cmd = m.textArea.Update(msg)
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
		
		content := strings.TrimSpace(string(updatedContent))
		if content != "" {
			err = database.UpdateNote(m.editingID, content)
			if err != nil {
				m.err = err
				return m, nil
			}
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
			m.textArea.Focus()
			return m, textarea.Blink
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
			titleStyle.Render("--- New Entry ---"),
			textAreaStyle.Render(m.textArea.View()),
			"(esc to cancel • ctrl+s to save)",
		)
	}

	// View for List Mode
	if len(m.notes) == 0 {
		return "No notes yet. Press 'n' to add one!\n\n(press q to quit)"
	}

	var s strings.Builder
	s.WriteString(titleStyle.Render("--- Your Notes ---") + "\n\n")

	for i, note := range m.notes {
		cursor := "  "
		displayContent := strings.ReplaceAll(note.Content, "\n", " ")
		displayContent = strings.ReplaceAll(displayContent, "\\n", " ")
		
		if len(displayContent) > 60 {
			displayContent = displayContent[:57] + "..."
		}

		if m.cursor == i {
			cursor = "> "
			s.WriteString(selectedStyle.Render(fmt.Sprintf("%s%s", cursor, displayContent)) + "\n")
		} else {
			s.WriteString(normalStyle.Render(fmt.Sprintf("%s%s", cursor, displayContent)) + "\n")
		}
	}

	selectedNote := m.notes[m.cursor]
	previewContent := strings.ReplaceAll(selectedNote.Content, "\\n", "\n")
	rendered, _ := glamour.Render(previewContent, "dark")
	s.WriteString("\n" + previewStyle.Render(rendered))

	// Status Bar
	help := "n: New • e: Edit • x: Delete • j/k: Nav • q: Quit"
	statusBar := statusBarStyle.Render(statusKey.Render(" JOTCLI ") + help)

	return s.String() + "\n" + statusBar
}