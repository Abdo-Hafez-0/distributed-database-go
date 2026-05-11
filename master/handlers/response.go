package handlers

import (
	"encoding/json"
	"net/http"

	"distributed-database-go/shared/types"
)

func WriteJSON(w http.ResponseWriter, status int, success bool, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(types.Response{
		Success: success,
		Message: message,
	})
}
