package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

func parseBody(r *http.Request, dst interface{}) error {
	if r.Body == nil {
		return errors.New("empty request body")
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return err
	}
	return nil
}

func methodNotAllowed(w http.ResponseWriter, expected string) {
	WriteJSON(w, http.StatusMethodNotAllowed, false, "method not allowed, expected "+expected)
}

func buildWhereClause(where map[string]interface{}) (string, []interface{}, error) {
	if len(where) == 0 {
		return "", nil, errors.New("where clause is required")
	}

	// Sort keys for deterministic query generation
	keys := make([]string, 0, len(where))
	for k := range where {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	args := make([]interface{}, 0, len(keys))

	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s = ?", key))
		args = append(args, where[key])
	}

	return strings.Join(parts, " AND "), args, nil
}
