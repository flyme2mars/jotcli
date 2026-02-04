package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type Note struct {
	ID        int
	Content   string
	Tag       string
	Priority  string
	CreatedAt time.Time
}

var DB *sql.DB

func InitDB() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dbPath := filepath.Join(home, ".jot.db")
	
	var errOpen error
	DB, errOpen = sql.Open("sqlite", dbPath)
	if errOpen != nil {
		return errOpen
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		tag TEXT,
		priority TEXT,
		created_at DATETIME
	);`

	_, err = DB.Exec(createTableSQL)
	return err
}

func AddNote(content, tag, priority string) error {
	query := `INSERT INTO notes (content, tag, priority, created_at) VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(query, content, tag, priority, time.Now())
	if err != nil {
		return fmt.Errorf("could not save note: %v", err)
	}
	return nil
}
