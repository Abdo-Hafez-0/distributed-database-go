package crud

import (
	"database/sql"
	"distributed-database-go/shared/utils"
	"fmt"
)

// QueryResult wraps the result of an INSERT / UPDATE / DELETE.
type QueryResult struct {
	LastInsertID int64
	RowsAffected int64
}

// InsertRecord inserts a single row into tableName.
// data maps column names to values; values are bound as parameters.
//
// Example:
//
//	InsertRecord(db, "users", map[string]interface{}{"name": "Alice", "age": 30})
func InsertRecord(db *sql.DB, tableName string, data map[string]interface{}) (QueryResult, error) {
	quotedTable, err := utils.QuoteTableName(tableName)
	if err != nil {
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
		"INSERT INTO %s (%s) VALUES (%s)",
		quotedTable, cols, placeholders,
	)

	res, err := db.Exec(query, args...)
	if err != nil {
		return QueryResult{}, fmt.Errorf("InsertRecord %q: %w", tableName, err)
	}

	lastID, _ := res.LastInsertId()
	rows, _ := res.RowsAffected()
	return QueryResult{LastInsertID: lastID, RowsAffected: rows}, nil
}
