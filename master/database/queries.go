package database

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

// QueryResult wraps the result of an INSERT / UPDATE / DELETE.
type QueryResult struct {
	LastInsertID int64
	RowsAffected int64
}

// ─────────────────────────────────────────────
// DDL — Data Definition
// ─────────────────────────────────────────────

// CreateDatabase creates a MySQL database if it does not already exist.
// Database names follow the same validation rules as table names.
func CreateDatabase(db *sql.DB, dbName string) error {
	if err := utils.ValidateIdentifier(dbName); err != nil {
		return fmt.Errorf("CreateDatabase: %w", err)
	}
	// Identifiers cannot be parameterised in DDL; we validated above.
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
// from a TableSchema.  Column definitions are validated before use.
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

// ─────────────────────────────────────────────
// DML — Data Manipulation
// ─────────────────────────────────────────────

// InsertRecord inserts a single row into tableName.
// data maps column names to values; values are bound as parameters.
//
// Example:
//
//	InsertRecord(db, "users", map[string]interface{}{"name": "Alice", "age": 30})
func InsertRecord(db *sql.DB, tableName string, data map[string]interface{}) (QueryResult, error) {
	if err := utils.ValidateIdentifier(tableName); err != nil {
		return QueryResult{}, fmt.Errorf("InsertRecord: %w", err)
	}
	if len(data) == 0 {
		return QueryResult{}, fmt.Errorf("InsertRecord %q: data map is empty", tableName)
	}

	cols, placeholders, args, err := utils.BuildInsertParts(data)
	if err != nil {
		return QueryResult{}, fmt.Errorf("InsertRecord %q: %w", tableName, err)
	}

	query := fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES (%s)",
		tableName,
		cols,
		placeholders,
	)

	res, err := db.Exec(query, args...)
	if err != nil {
		return QueryResult{}, fmt.Errorf("InsertRecord %q: %w", tableName, err)
	}

	lastID, _ := res.LastInsertId()
	rows, _ := res.RowsAffected()
	return QueryResult{LastInsertID: lastID, RowsAffected: rows}, nil
}

// UpdateRecord updates rows in tableName that match whereClause.
// data holds the columns to SET; whereArgs are bound to whereClause.
//
// Example:
//
//	UpdateRecord(db, "users", map[string]interface{}{"name": "Bob"}, "id = ?", 1)
func UpdateRecord(db *sql.DB, tableName string, data map[string]interface{}, whereClause string, whereArgs ...interface{}) (QueryResult, error) {
	if err := utils.ValidateIdentifier(tableName); err != nil {
		return QueryResult{}, fmt.Errorf("UpdateRecord: %w", err)
	}
	if len(data) == 0 {
		return QueryResult{}, fmt.Errorf("UpdateRecord %q: data map is empty", tableName)
	}
	if strings.TrimSpace(whereClause) == "" {
		return QueryResult{}, fmt.Errorf("UpdateRecord %q: whereClause must not be empty (full-table updates are not allowed)", tableName)
	}

	setClauses, setArgs, err := utils.BuildSetParts(data)
	if err != nil {
		return QueryResult{}, fmt.Errorf("UpdateRecord %q: %w", tableName, err)
	}

	query := fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", tableName, setClauses, whereClause)
	args := append(setArgs, whereArgs...)

	res, err := db.Exec(query, args...)
	if err != nil {
		return QueryResult{}, fmt.Errorf("UpdateRecord %q: %w", tableName, err)
	}

	lastID, _ := res.LastInsertId()
	rows, _ := res.RowsAffected()
	return QueryResult{LastInsertID: lastID, RowsAffected: rows}, nil
}

// DeleteRecord removes rows from tableName matching whereClause.
// whereClause must not be empty to prevent accidental full-table deletes.
//
// Example:
//
//	DeleteRecord(db, "users", "id = ?", 5)
func DeleteRecord(db *sql.DB, tableName string, whereClause string, whereArgs ...interface{}) (QueryResult, error) {
	if err := utils.ValidateIdentifier(tableName); err != nil {
		return QueryResult{}, fmt.Errorf("DeleteRecord: %w", err)
	}
	if strings.TrimSpace(whereClause) == "" {
		return QueryResult{}, fmt.Errorf("DeleteRecord %q: whereClause must not be empty", tableName)
	}

	query := fmt.Sprintf("DELETE FROM `%s` WHERE %s", tableName, whereClause)

	res, err := db.Exec(query, whereArgs...)
	if err != nil {
		return QueryResult{}, fmt.Errorf("DeleteRecord %q: %w", tableName, err)
	}

	lastID, _ := res.LastInsertId()
	rows, _ := res.RowsAffected()
	return QueryResult{LastInsertID: lastID, RowsAffected: rows}, nil
}

// SelectOptions configures a SELECT query.
type SelectOptions struct {
	Columns   []string // nil or empty → SELECT *
	Where     string   // optional WHERE clause (without the WHERE keyword)
	WhereArgs []interface{}
	OrderBy   string // e.g. "created_at DESC"
	Limit     int    // 0 means no LIMIT
	Offset    int
}

// SelectRecords executes a SELECT and returns rows as a slice of
// map[string]interface{}.  Column names are used as map keys.
//
// Example:
//
//	rows, err := SelectRecords(db, "users", SelectOptions{Where: "active = ?", WhereArgs: []interface{}{1}})
func SelectRecords(db *sql.DB, tableName string, opts SelectOptions) ([]map[string]interface{}, error) {
	if err := utils.ValidateIdentifier(tableName); err != nil {
		return nil, fmt.Errorf("SelectRecords: %w", err)
	}

	colStr := "*"
	if len(opts.Columns) > 0 {
		quoted := make([]string, 0, len(opts.Columns))
		for _, c := range opts.Columns {
			if err := utils.ValidateIdentifier(c); err != nil {
				return nil, fmt.Errorf("SelectRecords %q: invalid column %q: %w", tableName, c, err)
			}
			quoted = append(quoted, fmt.Sprintf("`%s`", c))
		}
		colStr = strings.Join(quoted, ", ")
	}

	query := fmt.Sprintf("SELECT %s FROM `%s`", colStr, tableName)

	if w := strings.TrimSpace(opts.Where); w != "" {
		query += " WHERE " + w
	}
	if ob := strings.TrimSpace(opts.OrderBy); ob != "" {
		query += " ORDER BY " + ob
	}
	if opts.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", opts.Limit)
		if opts.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", opts.Offset)
		}
	}

	rows, err := db.Query(query, opts.WhereArgs...)
	if err != nil {
		return nil, fmt.Errorf("SelectRecords %q: %w", tableName, err)
	}
	defer rows.Close()

	return utils.ScanRows(rows)
}
