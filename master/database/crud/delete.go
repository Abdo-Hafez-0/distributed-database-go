package crud

import (
	"database/sql"
	"distributed-database-go/shared/utils"
	"fmt"
	"strings"
)

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