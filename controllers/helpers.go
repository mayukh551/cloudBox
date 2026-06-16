package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

func respondWithJSON(w http.ResponseWriter, payload any, code int) {

	if http.StatusText(code) == "" {
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if payload == nil {
		json.NewEncoder(w).Encode(map[string]any{
			"success": true,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    payload,
	})
}

func respondWithError(w http.ResponseWriter, message any, code int, err error) {

	if message == nil || reflect.TypeOf(message).Kind() != reflect.String {
		fmt.Println("message not be nil and has to be string type")
		message = ""
	}

	fmt.Println("\n[Error Message] => ", (message))
	fmt.Println("[Raw Error] => ", err, "\n")

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
