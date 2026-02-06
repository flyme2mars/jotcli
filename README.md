# jotcli 📝

A minimalist, high-performance CLI tool for capturing thoughts and notes directly from your terminal. Built with Go, SQLite, and Bubble Tea.

## Features

- **🚀 Quick Capture**: Add notes instantly from the command line or within the app.
- **🎨 Markdown Preview**: Real-time rich text rendering for your notes using Glamour.
- **🛠 Interactive Dashboard**: A full-featured TUI to manage your thoughts without leaving the terminal.
- **📱 Responsive Design**: The list view automatically adapts to your terminal window size.
- **🔍 Full-Text Search**: Instantly find notes using the built-in search engine.
- **💾 Persistent Storage**: All data is saved securely in a local SQLite database (`~/.jot.db`).

## Installation

Ensure you have [Go](https://go.dev/doc/install) installed.

```bash
git clone https://github.com/flyme2mars/jotcli.git
cd jotcli
go install .
```

Alternatively, download the latest binary for your OS from the [Releases](https://github.com/flyme2mars/jotcli/releases) page.

## Usage

### Interactive Dashboard (Recommended)
Launch the main interface to browse, read, edit, and create notes.
```bash
jotcli view
```
**Shortcuts:**
- **↑/↓ or j/k**: Navigate notes
- **n**: Create a new note instantly
- **e**: Edit the selected note in your default editor ($EDITOR)
- **x or Backspace**: Delete the selected note
- **q or Ctrl+C**: Quit

### Command Line Interface

**Add a Note**
```bash
jotcli add "Check out the new #Go release" --tag dev --priority high
```

**Search Notes**
```bash
jotcli search "API"
```

**List & Filter**
```bash
jotcli list --tag work
```

**Edit by ID**
```bash
jotcli edit 5
```

## Tech Stack

- **Go**: High-performance systems language.
- **Cobra**: Industry-standard CLI framework.
- **SQLite**: Zero-config, reliable local database.
- **Bubble Tea**: Functional TUI framework for interactive elements.
- **Glamour**: Markdown rendering for the terminal.
- **Lip Gloss**: Modern terminal styling.

## License
MIT