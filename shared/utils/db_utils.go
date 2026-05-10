package utils

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

// ─────────────────────────────────────────────
// Validation
// ─────────────────────────────────────────────

// identifierRe allows letters, digits, and underscores only.
// This covers table names, column names, and database names.
var identifierRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{0,63}$`)

// allowedColumnTypes is a whitelist of SQL type keywords used in column definitions.
var allowedColumnTypes = []string{
	"INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT",
	"FLOAT", "DOUBLE", "DECIMAL", "NUMERIC",
	"VARCHAR", "CHAR", "TEXT", "MEDIUMTEXT", "LONGTEXT",
	"DATE", "DATETIME", "TIMESTAMP", "TIME", "YEAR",
	"BOOLEAN", "BOOL",
	"BLOB", "MEDIUMBLOB", "LONGBLOB",
	"JSON",
	"ENUM", "SET",
}

// ValidateIdentifier returns an error if name is not a safe SQL identifier.
// Safe means: starts with a letter or underscore, contains only alphanumeric
// characters and underscores, and is at most 64 characters long.
func ValidateIdentifier(name string) error {
	if name == "" {
		return fmt.Errorf("identifier must not be empty")
	}
	if !identifierRe.MatchString(name) {
		return fmt.Errorf("identifier %q contains invalid characters (only [A-Za-z0-9_] allowed, max 64 chars)", name)
	}
	return nil
}

// ValidateColumnDefinition does a lightweight check that a column definition
// string starts with a known SQL type keyword.  It does NOT parse the full
// SQL grammar — its purpose is to catch obvious injections or typos.
func ValidateColumnDefinition(def string) error {
	if def == "" {
		return fmt.Errorf("column definition must not be empty")
	}
	upper := strings.ToUpper(strings.TrimSpace(def))
	for _, t := range allowedColumnTypes {
		if strings.HasPrefix(upper, t) {
			return nil
		}
	}
	return fmt.Errorf("column definition %q does not start with a recognised SQL type", def)
}

// ─────────────────────────────────────────────
// Query fragment builders
// ─────────────────────────────────────────────

// BuildInsertParts converts a data map into the three parts of an INSERT:
//   - cols        → "`col1`, `col2`"
//   - placeholders → "?, ?"
//   - args        → []interface{}{val1, val2}
//
// Column names are validated before use.
func BuildInsertParts(data map[string]interface{}) (cols, placeholders string, args []interface{}, err error) {
	colNames := make([]string, 0, len(data))
	ph := make([]string, 0, len(data))
	vals := make([]interface{}, 0, len(data))

	for col, val := range data {
		if e := ValidateIdentifier(col); e != nil {
			return "", "", nil, fmt.Errorf("invalid column name %q: %w", col, e)
		}
		colNames = append(colNames, fmt.Sprintf("`%s`", col))
		ph = append(ph, "?")
		vals = append(vals, val)
	}

	return strings.Join(colNames, ", "), strings.Join(ph, ", "), vals, nil
}

// BuildSetParts converts a data map into the SET clause of an UPDATE:
//   - setClauses → "`col1` = ?, `col2` = ?"
//   - args       → []interface{}{val1, val2}
//
// Column names are validated before use.
func BuildSetParts(data map[string]interface{}) (setClauses string, args []interface{}, err error) {
	parts := make([]string, 0, len(data))
	vals := make([]interface{}, 0, len(data))

	for col, val := range data {
		if e := ValidateIdentifier(col); e != nil {
			return "", nil, fmt.Errorf("invalid column name %q: %w", col, e)
		}
		parts = append(parts, fmt.Sprintf("`%s` = ?", col))
		vals = append(vals, val)
	}

	return strings.Join(parts, ", "), vals, nil
}

// ─────────────────────────────────────────────
// Row scanning
// ─────────────────────────────────────────────

// ScanRows iterates over *sql.Rows and returns each row as a
// map[string]interface{} keyed by column name.
// The caller is responsible for closing rows AFTER this call returns.
func ScanRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("ScanRows: could not get columns: %w", err)
	}

	results := make([]map[string]interface{}, 0)

	for rows.Next() {
		// Allocate a fresh set of scan targets for each row.
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("ScanRows: scan error: %w", err)
		}

		row := make(map[string]interface{}, len(columns))
		for i, col := range columns {
			val := values[i]
			// MySQL driver returns []byte for strings; convert to string for usability.
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ScanRows: rows iteration error: %w", err)
	}

	return results, nil
}

// ─────────────────────────────────────────────
// String helpers
// ─────────────────────────────────────────────

// QuoteIdentifier wraps a validated identifier in backticks.
// It calls ValidateIdentifier first and returns an error on failure.
func QuoteIdentifier(name string) (string, error) {
	if err := ValidateIdentifier(name); err != nil {
		return "", err
	}
	return fmt.Sprintf("`%s`", name), nil
}

// PlaceholderList returns a comma-separated list of n question marks,
// useful for IN clauses.  E.g. PlaceholderList(3) → "?, ?, ?"
func PlaceholderList(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat("?, ", n-1) + "?"
}

// StringSliceToInterface converts []string to []interface{} for use as
// variadic query arguments.
func StringSliceToInterface(ss []string) []interface{} {
	out := make([]interface{}, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}