package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/flyme2mars/jotcli/internal/database"
	"github.com/spf13/cobra"
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
	
	
	
	func TestSearchIntegration(t *testing.T) {
	
		tempDB := "test_search.db"
	
		defer os.Remove(tempDB)
	
	
	
		var err error
	
		database.DB, err = database.OpenDB(tempDB)
	
		if err != nil {
	
			t.Fatalf("Failed to open test database: %v", err)
	
		}
	
		defer database.DB.Close()
	
	
	
		database.DB.Exec(`CREATE TABLE IF NOT EXISTS notes (
	
			id INTEGER PRIMARY KEY AUTOINCREMENT,
	
			content TEXT NOT NULL,
	
			tag TEXT,
	
			priority TEXT,
	
			created_at DATETIME
	
		);`)
	
	
	
			// 1. Add some notes
	
	
	
			database.AddNote("Apple pie recipe", "food", "low")
	
	
	
			database.AddNote("Banana bread", "food", "medium")
	
	
	
			database.AddNote("Buy a new computer", "work", "high")
	
	
	
		
	
	
	
			// Override the PersistentPreRun so it doesn't try to open ~/.jot.db
	
	
	
			rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {}
	
	
	
		
	
	
	
			// 2. Test searching for "Apple"
	
	
	
			buf := new(bytes.Buffer)
	
	
	
			rootCmd.SetOut(buf)
	
	
	
			rootCmd.SetErr(buf)
	
	
	
			rootCmd.SetArgs([]string{"search", "Apple"})
	
		
	
		err = rootCmd.Execute()
	
		if err != nil {
	
			t.Errorf("search command failed: %v", err)
	
		}
	
	
	
		output := buf.String()
	
		if !strings.Contains(output, "Apple pie recipe") {
	
			t.Errorf("Search failed to find 'Apple pie'. Output: %q", output)
	
		}
	
		if strings.Contains(output, "Banana bread") {
	
			t.Errorf("Search found 'Banana' when searching for 'Apple'. Output: %q", output)
	
		}
	
	
	
		// 3. Test searching for something that doesn't exist
	
		buf.Reset()
	
		rootCmd.SetArgs([]string{"search", "Zebra"})
	
		rootCmd.Execute()
	
		
	
		output = buf.String()
	
		if !strings.Contains(output, "No notes found matching 'Zebra'") {
	
			t.Errorf("Search should have returned no results. Output: %q", output)
	
		}
	
	}
	
	