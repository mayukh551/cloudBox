package controllers

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, payload any, code int) {

	if http.StatusText(code) == "" {
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    payload,
	})
}

func respondWithError(w http.ResponseWriter, message any, code int) {

	if http.StatusText(code) == "" {
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error":   message,
	})
}
