package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

func respondJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"data": body,
		"meta": map[string]string{
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   VersionString(),
		},
	})
}
