package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/flyme2mars/jotcli/internal/database"
)

func TestAddAndListIntegration(t *testing.T) {
	// 1. Setup: Use a temporary database
	tempDB := "test_integration.db"
	defer os.Remove(tempDB)

	var err error
	database.DB, err = database.OpenDB(tempDB)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer database.DB.Close()

	// Initialize tables
	createTableSQL := `CREATE TABLE IF NOT EXISTS notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		tag TEXT,
		priority TEXT,
		created_at DATETIME
	);`
	database.DB.Exec(createTableSQL)

	// 2. Test "add" command
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	
	// Reset flags to avoid contamination from previous tests
	tag = ""
	priority = "low"
	
	rootCmd.SetArgs([]string{"add", "Integration Test Note", "--tag", "test"})
	
	err = rootCmd.Execute()
	if err != nil {
		t.Errorf("add command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "✅ Note saved") {
		t.Errorf("Expected success message in output, but didn't find it. Output: %q", output)
	}

	// 3. Test "list" command
	buf.Reset()
	rootCmd.SetArgs([]string{"list"})
	
	err = rootCmd.Execute()
	if err != nil {
		t.Errorf("list command failed: %v", err)
	}

	output = buf.String()
	if !strings.Contains(output, "Integration Test Note") {
		t.Errorf("List output missing added note. Output: %q", output)
	}
	
	if !strings.Contains(output, "test") {
		t.Errorf("List output missing tag. Output: %q", output)
	}
}