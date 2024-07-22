package database

import (
	"os"
	"testing"
)

func TestNewDB(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := tmpDir + "/database.json"

	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	data, err := os.ReadFile(dbPath)
	if err != nil {
		t.Fatalf("Expected no error: got %v", err)
	}

	expectedContent := `{"chirps": {}}`
	if string(data) != expectedContent {
		t.Fatalf("Expected content %v: got %v", expectedContent, string(data))
	}

	if db.path != dbPath {
		t.Fatalf("Expected path %v: got %v", dbPath, db.path)
	}

	if db.mux == nil {
		t.Fatal("Expected mux to be non-nil")
	}
}
