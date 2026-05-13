package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"distributed-database-go/master/database"
	"distributed-database-go/shared/types"
	"distributed-database-go/shared/utils"
)

var db *sql.DB

func SetDatabase(databaseConn *sql.DB) {
	db = databaseConn
}

func HandleReplication(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req types.ReplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Replication Request -> Operation: %s | Database: %s | Table: %s\n", req.Operation, req.Database, req.Table)

	if db == nil {
		http.Error(w, "database connection is not initialized", http.StatusInternalServerError)
		return
	}

	var err error
	switch req.Operation {
	case "insert":
		err = handleInsert(req)
	case "update":
		err = handleUpdate(req)
	case "delete":
		err = handleDelete(req)
	default:
		http.Error(w, "Unknown operation", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("Replication error: %v", err)
		http.Error(w, "Replication failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := types.ReplicationResponse{
		Status:  "success",
		Message: "Replication completed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HandleSelect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if db == nil {
		http.Error(w, "database connection is not initialized", http.StatusInternalServerError)
		return
	}

	params := r.URL.Query()
	table := params.Get("table")
	if table == "" {
		http.Error(w, "table query parameter is required", http.StatusBadRequest)
		return
	}

	databaseName := params.Get("database")
	targetTable := qualifiedTableName(databaseName, table)
	limit := parseInt(params.Get("limit"), 100)
	offset := parseInt(params.Get("offset"), 0)

	results, err := database.SelectRecords(db, targetTable, database.SelectOptions{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		http.Error(w, "Select failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    results,
	})
}

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if db == nil {
		http.Error(w, "database connection is not initialized", http.StatusInternalServerError)
		return
	}

	params := r.URL.Query()
	table := params.Get("table")
	if table == "" {
		http.Error(w, "table query parameter is required", http.StatusBadRequest)
		return
	}

	queryText := params.Get("q")
	if queryText == "" {
		http.Error(w, "q query parameter is required", http.StatusBadRequest)
		return
	}

	databaseName := params.Get("database")
	targetTable := qualifiedTableName(databaseName, table)
	limit := parseInt(params.Get("limit"), 100)
	offset := parseInt(params.Get("offset"), 0)
	column := params.Get("column")

	whereClause, whereArgs, err := buildSearchWhere(targetTable, column, queryText)
	if err != nil {
		http.Error(w, "Search failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	results, err := database.SelectRecords(db, targetTable, database.SelectOptions{
		Where:     whereClause,
		WhereArgs: whereArgs,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		http.Error(w, "Search failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    results,
	})
}

func handleInsert(req types.ReplicationRequest) error {
	if len(req.Data) == 0 {
		return errors.New("insert payload is empty")
	}

	tableName := qualifiedTableName(req.Database, req.Table)
	if _, err := database.InsertRecord(db, tableName, req.Data); err != nil {
		return fmt.Errorf("insert replication: %w", err)
	}

	log.Printf("Insert replicated to %s", tableName)
	return nil
}

func handleUpdate(req types.ReplicationRequest) error {
	rawWhere, ok := req.Data["where"]
	if !ok {
		return errors.New("update replication requires where clause")
	}

	rawData, ok := req.Data["data"]
	if !ok {
		return errors.New("update replication requires data payload")
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		return errors.New("update data must be an object")
	}

	whereClause, whereArgs, err := buildWhereClause(rawWhere)
	if err != nil {
		return err
	}

	tableName := qualifiedTableName(req.Database, req.Table)
	if _, err := database.UpdateRecord(db, tableName, data, whereClause, whereArgs...); err != nil {
		return fmt.Errorf("update replication: %w", err)
	}

	log.Printf("Update replicated to %s", tableName)
	return nil
}

func handleDelete(req types.ReplicationRequest) error {
	rawWhere, ok := req.Data["where"]
	if !ok {
		return errors.New("delete replication requires where clause")
	}

	whereClause, whereArgs, err := buildWhereClause(rawWhere)
	if err != nil {
		return err
	}

	tableName := qualifiedTableName(req.Database, req.Table)
	if _, err := database.DeleteRecord(db, tableName, whereClause, whereArgs...); err != nil {
		return fmt.Errorf("delete replication: %w", err)
	}

	log.Printf("Delete replicated from %s", tableName)
	return nil
}

func qualifiedTableName(databaseName, table string) string {
	if databaseName != "" {
		return fmt.Sprintf("%s.%s", databaseName, table)
	}
	return table
}

func buildWhereClause(raw interface{}) (string, []interface{}, error) {
	whereMap, ok := raw.(map[string]interface{})
	if !ok {
		return "", nil, errors.New("where clause must be an object")
	}

	if len(whereMap) == 0 {
		return "", nil, errors.New("where clause must not be empty")
	}

	keys := make([]string, 0, len(whereMap))
	for k := range whereMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	args := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		if err := utils.ValidateIdentifier(key); err != nil {
			return "", nil, fmt.Errorf("invalid where field %q: %w", key, err)
		}
		parts = append(parts, fmt.Sprintf("`%s` = ?", key))
		args = append(args, whereMap[key])
	}

	return strings.Join(parts, " AND "), args, nil
}

func buildSearchWhere(tableName, column, queryText string) (string, []interface{}, error) {
	if column != "" {
		if err := utils.ValidateIdentifier(column); err != nil {
			return "", nil, fmt.Errorf("invalid search column: %w", err)
		}
		return fmt.Sprintf("`%s` LIKE ?", column), []interface{}{"%" + queryText + "%"}, nil
	}

	columns, err := loadTableColumns(tableName)
	if err != nil {
		return "", nil, err
	}

	if len(columns) == 0 {
		return "", nil, errors.New("no columns found for search")
	}

	quotedColumns := make([]string, 0, len(columns))
	for _, c := range columns {
		quotedColumns = append(quotedColumns, fmt.Sprintf("`%s`", c))
	}

	return fmt.Sprintf("CONCAT_WS(' ', %s) LIKE ?", strings.Join(quotedColumns, ", ")), []interface{}{"%" + queryText + "%"}, nil
}

func loadTableColumns(tableName string) ([]string, error) {
	quotedTable, err := utils.QuoteTableName(tableName)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SHOW COLUMNS FROM " + quotedTable)
	if err != nil {
		return nil, fmt.Errorf("could not load table columns: %w", err)
	}
	defer rows.Close()

	columns := make([]string, 0)
	for rows.Next() {
		var field string
		var colType, nullable, key, defaultValue, extra sql.NullString
		if err := rows.Scan(&field, &colType, &nullable, &key, &defaultValue, &extra); err != nil {
			return nil, fmt.Errorf("could not scan columns metadata: %w", err)
		}
		columns = append(columns, field)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading table columns: %w", err)
	}

	return columns, nil
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return defaultValue
	}
	return parsed
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
