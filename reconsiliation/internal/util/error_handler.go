package util

import (
	"encoding/json"
	"net/http"
)

func HandleError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"message": msg,
	})
}
