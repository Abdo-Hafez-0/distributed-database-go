// Package database is the public entry point for all database operations.
// All CRUD and DDL functions live under the crud sub-package and are
// re-exported here so callers import only "distributed-database-go/master/database".
package database

import (
	"database/sql"
	"distributed-database-go/master/database/crud"
)

// ── Type aliases ─────────────────────────────────────────────────────────────

type ColumnDefinition = crud.ColumnDefinition
type TableSchema = crud.TableSchema
type QueryResult = crud.QueryResult
type SelectOptions = crud.SelectOptions

// ── DDL ──────────────────────────────────────────────────────────────────────

func CreateDatabase(db *sql.DB, dbName string) error {
	return crud.CreateDatabase(db, dbName)
}

func CreateTable(db *sql.DB, schema TableSchema) error {
	return crud.CreateTable(db, schema)
}

// ── DML ──────────────────────────────────────────────────────────────────────

func InsertRecord(db *sql.DB, tableName string, data map[string]interface{}) (QueryResult, error) {
	return crud.InsertRecord(db, tableName, data)
}

func UpdateRecord(db *sql.DB, tableName string, data map[string]interface{}, whereClause string, whereArgs ...interface{}) (QueryResult, error) {
	return crud.UpdateRecord(db, tableName, data, whereClause, whereArgs...)
}

func DeleteRecord(db *sql.DB, tableName string, whereClause string, whereArgs ...interface{}) (QueryResult, error) {
	return crud.DeleteRecord(db, tableName, whereClause, whereArgs...)
}

func SelectRecords(db *sql.DB, tableName string, opts SelectOptions) ([]map[string]interface{}, error) {
	return crud.SelectRecords(db, tableName, opts)
}