package database_test

import (
	"distributed-database-go/master/database"
	"distributed-database-go/shared/utils"
	"os"
	"testing"
)

// ─────────────────────────────────────────────
// Test helpers
// ─────────────────────────────────────────────

// setupTestDB opens a real connection using env vars.
// Tests are skipped when DB_HOST is not set so CI without a MySQL
// instance does not fail.
func setupTestDB(t *testing.T) {
	t.Helper()
	if os.Getenv("DB_HOST") == "" {
		t.Skip("DB_HOST not set — skipping integration tests")
	}
	database.MustConnect()
}

// testTable is a throwaway table name used across tests.
const testTable = "db_layer_test"

// dropTestTable cleans up after each test run.
func dropTestTable(t *testing.T) {
	t.Helper()
	_, err := database.DB.Exec("DROP TABLE IF EXISTS `" + testTable + "`")
	if err != nil {
		t.Fatalf("cleanup: could not drop %s: %v", testTable, err)
	}
}

// ─────────────────────────────────────────────
// Tests
// ─────────────────────────────────────────────

func TestConnect(t *testing.T) {
	setupTestDB(t)
	if err := database.DB.Ping(); err != nil {
		t.Fatalf("Ping after Connect failed: %v", err)
	}
}

func TestCreateTable(t *testing.T) {
	setupTestDB(t)
	defer dropTestTable(t)

	schema := database.TableSchema{
		TableName: testTable,
		Columns: []database.ColumnDefinition{
			{Name: "id",    Definition: "INT NOT NULL AUTO_INCREMENT PRIMARY KEY"},
			{Name: "name",  Definition: "VARCHAR(255) NOT NULL"},
			{Name: "email", Definition: "VARCHAR(255)"},
			{Name: "score", Definition: "DECIMAL(10,2) DEFAULT 0.00"},
		},
	}

	if err := database.CreateTable(database.DB, schema); err != nil {
		t.Fatalf("CreateTable: %v", err)
	}

	// Running again must be idempotent (IF NOT EXISTS).
	if err := database.CreateTable(database.DB, schema); err != nil {
		t.Fatalf("CreateTable (idempotent run): %v", err)
	}
}

func TestInsertRecord(t *testing.T) {
	setupTestDB(t)
	defer dropTestTable(t)

	// Ensure table exists.
	_ = database.CreateTable(database.DB, database.TableSchema{
		TableName: testTable,
		Columns: []database.ColumnDefinition{
			{Name: "id",    Definition: "INT NOT NULL AUTO_INCREMENT PRIMARY KEY"},
			{Name: "name",  Definition: "VARCHAR(255) NOT NULL"},
			{Name: "email", Definition: "VARCHAR(255)"},
		},
	})

	res, err := database.InsertRecord(database.DB, testTable, map[string]interface{}{
		"name":  "Alice",
		"email": "alice@example.com",
	})
	if err != nil {
		t.Fatalf("InsertRecord: %v", err)
	}
	if res.LastInsertID == 0 {
		t.Error("expected non-zero LastInsertID")
	}
	if res.RowsAffected != 1 {
		t.Errorf("expected 1 row affected, got %d", res.RowsAffected)
	}
}

func TestUpdateRecord(t *testing.T) {
	setupTestDB(t)
	defer dropTestTable(t)

	_ = database.CreateTable(database.DB, database.TableSchema{
		TableName: testTable,
		Columns: []database.ColumnDefinition{
			{Name: "id",   Definition: "INT NOT NULL AUTO_INCREMENT PRIMARY KEY"},
			{Name: "name", Definition: "VARCHAR(255) NOT NULL"},
		},
	})
	ins, _ := database.InsertRecord(database.DB, testTable, map[string]interface{}{"name": "Bob"})

	res, err := database.UpdateRecord(
		database.DB, testTable,
		map[string]interface{}{"name": "Robert"},
		"id = ?", ins.LastInsertID,
	)
	if err != nil {
		t.Fatalf("UpdateRecord: %v", err)
	}
	if res.RowsAffected != 1 {
		t.Errorf("expected 1 row affected, got %d", res.RowsAffected)
	}
}

func TestDeleteRecord(t *testing.T) {
	setupTestDB(t)
	defer dropTestTable(t)

	_ = database.CreateTable(database.DB, database.TableSchema{
		TableName: testTable,
		Columns: []database.ColumnDefinition{
			{Name: "id",   Definition: "INT NOT NULL AUTO_INCREMENT PRIMARY KEY"},
			{Name: "name", Definition: "VARCHAR(255) NOT NULL"},
		},
	})
	ins, _ := database.InsertRecord(database.DB, testTable, map[string]interface{}{"name": "Charlie"})

	res, err := database.DeleteRecord(database.DB, testTable, "id = ?", ins.LastInsertID)
	if err != nil {
		t.Fatalf("DeleteRecord: %v", err)
	}
	if res.RowsAffected != 1 {
		t.Errorf("expected 1 row affected, got %d", res.RowsAffected)
	}
}

func TestSelectRecords(t *testing.T) {
	setupTestDB(t)
	defer dropTestTable(t)

	_ = database.CreateTable(database.DB, database.TableSchema{
		TableName: testTable,
		Columns: []database.ColumnDefinition{
			{Name: "id",   Definition: "INT NOT NULL AUTO_INCREMENT PRIMARY KEY"},
			{Name: "name", Definition: "VARCHAR(255) NOT NULL"},
			{Name: "active", Definition: "BOOLEAN NOT NULL DEFAULT TRUE"},
		},
	})

	names := []string{"Dave", "Eve", "Frank"}
	for _, n := range names {
		_, _ = database.InsertRecord(database.DB, testTable, map[string]interface{}{"name": n, "active": true})
	}
	_, _ = database.InsertRecord(database.DB, testTable, map[string]interface{}{"name": "Inactive", "active": false})

	rows, err := database.SelectRecords(database.DB, testTable, database.SelectOptions{
		Where:     "active = ?",
		WhereArgs: []interface{}{1},
		OrderBy:   "id ASC",
	})
	if err != nil {
		t.Fatalf("SelectRecords: %v", err)
	}
	if len(rows) != len(names) {
		t.Errorf("expected %d rows, got %d", len(names), len(rows))
	}
}

func TestSelectRecordsWithLimit(t *testing.T) {
	setupTestDB(t)
	defer dropTestTable(t)

	_ = database.CreateTable(database.DB, database.TableSchema{
		TableName: testTable,
		Columns: []database.ColumnDefinition{
			{Name: "id",   Definition: "INT NOT NULL AUTO_INCREMENT PRIMARY KEY"},
			{Name: "name", Definition: "VARCHAR(255) NOT NULL"},
		},
	})
	for i := 0; i < 10; i++ {
		_, _ = database.InsertRecord(database.DB, testTable, map[string]interface{}{"name": "User"})
	}

	rows, err := database.SelectRecords(database.DB, testTable, database.SelectOptions{Limit: 3})
	if err != nil {
		t.Fatalf("SelectRecords (limit): %v", err)
	}
	if len(rows) != 3 {
		t.Errorf("expected 3 rows, got %d", len(rows))
	}
}

// ─────────────────────────────────────────────
// Validation unit tests (no DB needed)
// ─────────────────────────────────────────────

func TestValidateIdentifier(t *testing.T) {
	cases := []struct {
		input string
		valid bool
	}{
		{"users", true},
		{"order_items", true},
		{"_private", true},
		{"", false},
		{"123abc", false},
		{"drop table", false},
		{"users; DROP TABLE users", false},
		{"a very long identifier that exceeds sixty four characters in total here", false},
	}

	for _, tc := range cases {
		err := utils.ValidateIdentifier(tc.input)
		if tc.valid && err != nil {
			t.Errorf("ValidateIdentifier(%q): expected valid, got error: %v", tc.input, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("ValidateIdentifier(%q): expected error, got nil", tc.input)
		}
	}
}

func TestValidateColumnDefinition(t *testing.T) {
	cases := []struct {
		input string
		valid bool
	}{
		{"INT NOT NULL AUTO_INCREMENT", true},
		{"VARCHAR(255) NOT NULL", true},
		{"DECIMAL(10,2) DEFAULT 0.00", true},
		{"BOOLEAN DEFAULT TRUE", true},
		{"", false},
		{"UNKNOWN_TYPE(255)", false},
	}

	for _, tc := range cases {
		err := utils.ValidateColumnDefinition(tc.input)
		if tc.valid && err != nil {
			t.Errorf("ValidateColumnDefinition(%q): expected valid, got: %v", tc.input, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("ValidateColumnDefinition(%q): expected error, got nil", tc.input)
		}
	}
}

func TestPlaceholderList(t *testing.T) {
	cases := map[int]string{
		0: "",
		1: "?",
		3: "?, ?, ?",
	}
	for n, want := range cases {
		got := utils.PlaceholderList(n)
		if got != want {
			t.Errorf("PlaceholderList(%d) = %q, want %q", n, got, want)
		}
	}
}