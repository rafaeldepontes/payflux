package util

import (
	"encoding/json"
	"net/http"
)

func HandleError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": msg,
	})
}
