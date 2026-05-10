package crud

import (
	"database/sql"
	"distributed-database-go/shared/utils"
	"fmt"
	"strings"
)

// ColumnDefinition describes a single column in a CREATE TABLE statement.
type ColumnDefinition struct {
	Name       string // e.g. "id"
	Definition string // e.g. "INT NOT NULL AUTO_INCREMENT PRIMARY KEY"
}

// TableSchema is the payload used by CreateTable.
type TableSchema struct {
	TableName string
	Columns   []ColumnDefinition
}

// CreateDatabase creates a MySQL database if it does not already exist.
func CreateDatabase(db *sql.DB, dbName string) error {
	if err := utils.ValidateIdentifier(dbName); err != nil {
		return fmt.Errorf("CreateDatabase: %w", err)
	}

	query := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		dbName,
	)

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("CreateDatabase %q: %w", dbName, err)
	}
	return nil
}

// CreateTable builds and executes a CREATE TABLE IF NOT EXISTS statement
// from a TableSchema. Column definitions are validated before use.
func CreateTable(db *sql.DB, schema TableSchema) error {
	if err := utils.ValidateIdentifier(schema.TableName); err != nil {
		return fmt.Errorf("CreateTable: invalid table name: %w", err)
	}
	if len(schema.Columns) == 0 {
		return fmt.Errorf("CreateTable %q: at least one column is required", schema.TableName)
	}

	colDefs := make([]string, 0, len(schema.Columns))
	for _, col := range schema.Columns {
		if err := utils.ValidateIdentifier(col.Name); err != nil {
			return fmt.Errorf("CreateTable %q: invalid column name %q: %w", schema.TableName, col.Name, err)
		}
		if err := utils.ValidateColumnDefinition(col.Definition); err != nil {
			return fmt.Errorf("CreateTable %q: invalid definition for column %q: %w", schema.TableName, col.Name, err)
		}
		colDefs = append(colDefs, fmt.Sprintf("  `%s` %s", col.Name, col.Definition))
	}

	query := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS `%s` (\n%s\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4",
		schema.TableName,
		strings.Join(colDefs, ",\n"),
	)

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("CreateTable %q: %w", schema.TableName, err)
	}
	return nil
}