# jotcli 📝

A minimalist, high-performance CLI tool for capturing thoughts and notes directly from your terminal. Built with Go, SQLite, and Bubble Tea.

## Features

- **Quick Capture**: Add notes instantly with tags and priority.
- **Persistent Storage**: All notes are saved in a local SQLite database (`~/.jot.db`).
- **Pretty Listing**: View your notes in a clean, formatted table.
- **Interactive TUI**: Browse and delete notes using a keyboard-driven interface.

## Installation

Ensure you have [Go](https://go.dev/doc/install) installed.

```bash
git clone https://github.com/flyme2mars/jotcli.git
cd jotcli
go install .
```

## Usage

### Add a Note
```bash
jotcli add "Fix the API bug" --tag work --priority high
```

### List Notes
```bash
# View all notes
jotcli list

# Filter by tag
jotcli list --tag work
```

### Interactive Mode
Launch the interactive viewer to browse and manage your notes.
```bash
jotcli view
```
- **↑/↓ or j/k**: Navigate
- **x or Backspace**: Delete selected note
- **q or Ctrl+C**: Quit

## Tech Stack

- **Go**: The core language.
- **Cobra**: CLI framework for commands and flags.
- **SQLite**: Reliable local data persistence.
- **Bubble Tea**: TUI framework for the interactive view.
- **Lip Gloss**: Terminal styling and layouts.

## License
MIT
