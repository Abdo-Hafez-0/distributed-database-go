package crud

import (
	"database/sql"
	"distributed-database-go/shared/utils"
	"fmt"
	"strings"
)

// UpdateRecord updates rows in tableName that match whereClause.
// data holds the columns to SET; whereArgs are bound to whereClause.
//
// Example:
//
//	UpdateRecord(db, "users", map[string]interface{}{"name": "Bob"}, "id = ?", 1)
func UpdateRecord(db *sql.DB, tableName string, data map[string]interface{}, whereClause string, whereArgs ...interface{}) (QueryResult, error) {
	quotedTable, err := utils.QuoteTableName(tableName)
	if err != nil {
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

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", quotedTable, setClauses, whereClause)
	args := append(setArgs, whereArgs...)

	res, err := db.Exec(query, args...)
	if err != nil {
		return QueryResult{}, fmt.Errorf("UpdateRecord %q: %w", tableName, err)
	}

	lastID, _ := res.LastInsertId()
	rows, _ := res.RowsAffected()
	return QueryResult{LastInsertID: lastID, RowsAffected: rows}, nil
}
