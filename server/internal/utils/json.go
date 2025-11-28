package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		slog.Error("Error writing JSON response", "error", err)
		//WriteJSONMessage(w, http.StatusInternalServerError, "the server encountered a problem")
	}
}

type ResponseMessage struct {
	Message string `json:"error"`
}

func WriteJSONMessage(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, &ResponseMessage{Message: message})
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}
