package database

import (
	"os"
	"testing"
)

func TestAddAndGetNotes(t *testing.T) {
	// 1. Setup: Use a temporary database for testing
	tempDB := "test_jot.db"
	// Ensure cleanup after test
	defer os.Remove(tempDB)

	// Mock the DB initialization logic for test
	var err error
	DB, err = OpenDB(tempDB)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer DB.Close()

	// Ensure table exists
	createTableSQL := `CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		tag TEXT,
		priority TEXT,
		created_at DATETIME
	);`
	DB.Exec(createTableSQL)

	// 2. Test Cases (Table-Driven)
	tests := []struct {
		name     string
		content  string
		tag      string
		priority string
	}{
		{"Simple Note", "Hello World", "test", "low"},
		{"Markdown Note", "# Title\n- Item", "work", "high"},
		{"Empty Tag", "No tag here", "", "medium"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test AddNote
			err := AddNote(tt.content, tt.tag, tt.priority)
			if err != nil {
				t.Errorf("AddNote() error = %v", err)
			}

			// Test retrieval
			// We use empty string to get all notes for verification
			notes, err := GetNotes("")
			if err != nil {
				t.Errorf("GetNotes() error = %v", err)
			}

			found := false
			for _, n := range notes {
				if n.Content == tt.content {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("GetNotes() did not return the added note: %s", tt.content)
			}
		})
	}
}

func TestDeleteNote(t *testing.T) {
	tempDB := "test_delete.db"
	defer os.Remove(tempDB)

	var err error
	DB, err = OpenDB(tempDB)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer DB.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		tag TEXT,
		priority TEXT,
		created_at DATETIME
	);`
	DB.Exec(createTableSQL)

	// Add a note to delete
	err = AddNote("Delete me", "trash", "low")
	if err != nil {
		t.Fatalf("Failed to add note: %v", err)
	}

	notes, _ := GetNotes("")
	if len(notes) != 1 {
		t.Fatalf("Expected 1 note, got %d", len(notes))
	}

	// Delete it
	err = DeleteNote(notes[0].ID)
	if err != nil {
		t.Errorf("DeleteNote() error = %v", err)
	}

	// Verify it's gone
	notes, _ = GetNotes("")
	if len(notes) != 0 {
		t.Errorf("Note was not deleted, still have %d notes", len(notes))
	}
}