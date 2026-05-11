package handlers

import (
	"net/http"

	"distributed-database-go/master/database"
	"distributed-database-go/master/replication"
	"distributed-database-go/shared/types"
)

func DropDatabaseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		methodNotAllowed(w, "DELETE")
		return
	}

	var req types.DropDatabaseRequest
	if err := parseBody(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, false, "invalid request body")
		return
	}

	if req.Database == "" {
		WriteJSON(w, http.StatusBadRequest, false, "database is required")
		return
	}

	if database.DB == nil {
		WriteJSON(w, http.StatusInternalServerError, false, "database connection is not initialized")
		return
	}

	if err := database.DropDatabase(database.DB, req.Database); err != nil {
		WriteJSON(w, http.StatusInternalServerError, false, err.Error())
		return
	}

	repReq := types.ReplicationRequest{
		Operation: "drop_database",
		Database:  req.Database,
	}

	if err := replication.Replicate(repReq); err != nil {
		WriteJSON(w, http.StatusInternalServerError, false, "database dropped locally but replication failed: "+err.Error())
		return
	}

	WriteJSON(w, http.StatusOK, true, "database dropped successfully")
}
