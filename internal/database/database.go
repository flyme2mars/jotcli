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
	
	DB, err = OpenDB(dbPath)
	if err != nil {
		return err
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

func OpenDB(path string) (*sql.DB, error) {
	return sql.Open("sqlite", path)
}

func AddNote(content, tag, priority string) error {
	query := `INSERT INTO notes (content, tag, priority, created_at) VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(query, content, tag, priority, time.Now())
	if err != nil {
		return fmt.Errorf("could not save note: %v", err)
	}
	return nil
}

func GetNotes(tagFilter string) ([]Note, error) {
	var rows *sql.Rows
	var err error

	if tagFilter != "" {
		query := `SELECT id, content, tag, priority, created_at FROM notes WHERE tag = ? ORDER BY created_at DESC`
		rows, err = DB.Query(query, tagFilter)
	} else {
		query := `SELECT id, content, tag, priority, created_at FROM notes ORDER BY created_at DESC`
		rows, err = DB.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		err := rows.Scan(&n.ID, &n.Content, &n.Tag, &n.Priority, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func GetNotesBySearch(query string) ([]Note, error) {
	sqlQuery := `SELECT id, content, tag, priority, created_at FROM notes WHERE content LIKE ? ORDER BY created_at DESC`
	rows, err := DB.Query(sqlQuery, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var n Note
		err := rows.Scan(&n.ID, &n.Content, &n.Tag, &n.Priority, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func GetNoteByID(id int) (*Note, error) {
	query := `SELECT id, content, tag, priority, created_at FROM notes WHERE id = ?`
	row := DB.QueryRow(query, id)

	var n Note
	err := row.Scan(&n.ID, &n.Content, &n.Tag, &n.Priority, &n.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &n, nil
}

func UpdateNote(id int, content string) error {
	query := `UPDATE notes SET content = ? WHERE id = ?`
	_, err := DB.Exec(query, content, id)
	return err
}

func DeleteNote(id int) error {
	query := `DELETE FROM notes WHERE id = ?`
	_, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not delete note: %v", err)
	}
	return nil
}
