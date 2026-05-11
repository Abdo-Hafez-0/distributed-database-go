package handlers

import (
	"net/http"

	"distributed-database-go/master/database"
	"distributed-database-go/master/replication"
	"distributed-database-go/shared/types"
)

func CreateTableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}

	var req types.CreateTableRequest
	if err := parseBody(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, false, "invalid request body")
		return
	}

	if req.Database == "" {
		WriteJSON(w, http.StatusBadRequest, false, "database is required")
		return
	}

	if req.Table == "" {
		WriteJSON(w, http.StatusBadRequest, false, "table is required")
		return
	}

	if len(req.Columns) == 0 {
		WriteJSON(w, http.StatusBadRequest, false, "columns are required")
		return
	}

	if database.DB == nil {
		WriteJSON(w, http.StatusInternalServerError, false, "database connection is not initialized")
		return
	}

	// Convert map[string]string to []ColumnDefinition
	columns := make([]database.ColumnDefinition, 0, len(req.Columns))
	for name, definition := range req.Columns {
		columns = append(columns, database.ColumnDefinition{
			Name:       name,
			Definition: definition,
		})
	}

	schema := database.TableSchema{
		TableName: req.Table,
		Columns:   columns,
	}

	if err := database.CreateTable(database.DB, schema); err != nil {
		WriteJSON(w, http.StatusInternalServerError, false, err.Error())
		return
	}

	repReq := types.ReplicationRequest{
		Operation: "create_table",
		Database:  req.Database,
		Table:     req.Table,
		Data:      map[string]interface{}{"columns": req.Columns},
	}

	if err := replication.Replicate(repReq); err != nil {
		WriteJSON(w, http.StatusInternalServerError, false, "table created locally but replication failed: "+err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, true, "table created successfully")
}
