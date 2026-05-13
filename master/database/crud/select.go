package crud

import (
	"database/sql"
	"distributed-database-go/shared/utils"
	"fmt"
	"strings"
)

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
// map[string]interface{}. Column names are used as map keys.
//
// Example:
//
//	SelectRecords(db, "users", SelectOptions{Where: "active = ?", WhereArgs: []interface{}{1}})
func SelectRecords(db *sql.DB, tableName string, opts SelectOptions) ([]map[string]interface{}, error) {
	quotedTable, err := utils.QuoteTableName(tableName)
	if err != nil {
		return nil, fmt.Errorf("SelectRecords: %w", err)
	}

	colStr, err := buildColumnList(opts.Columns, tableName)
	if err != nil {
		return nil, err
	}

	query := buildSelectQuery(quotedTable, colStr, opts)

	rows, err := db.Query(query, opts.WhereArgs...)
	if err != nil {
		return nil, fmt.Errorf("SelectRecords %q: %w", tableName, err)
	}
	defer rows.Close()

	return utils.ScanRows(rows)
}

// buildColumnList validates and quotes the requested columns.
// Returns "*" when no columns are specified.
func buildColumnList(columns []string, tableName string) (string, error) {
	if len(columns) == 0 {
		return "*", nil
	}

	quoted := make([]string, 0, len(columns))
	for _, c := range columns {
		if err := utils.ValidateIdentifier(c); err != nil {
			return "", fmt.Errorf("SelectRecords %q: invalid column %q: %w", tableName, c, err)
		}
		quoted = append(quoted, fmt.Sprintf("`%s`", c))
	}
	return strings.Join(quoted, ", "), nil
}

// buildSelectQuery assembles the full SELECT statement from its parts.
func buildSelectQuery(tableName, colStr string, opts SelectOptions) string {
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
	return query
}
