package util

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	tests := []struct {
		name           string
		msg            string
		code           int
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "Bad Request",
			msg:            "invalid input",
			code:           http.StatusBadRequest,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"message": "invalid input"},
		},
		{
			name:           "Internal Error",
			msg:            "something went wrong",
			code:           http.StatusInternalServerError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"message": "something went wrong"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			HandleError(w, tt.msg, tt.code)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var body map[string]string
			json.NewDecoder(w.Body).Decode(&body)
			assert.Equal(t, tt.expectedBody, body)
		})
	}
}
