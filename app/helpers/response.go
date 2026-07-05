package helper

import (
	"encoding/json"
	"net/http"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type APIResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func DecodeJSON(r *http.Request, target any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(target)
}

func WriteJSON(w http.ResponseWriter, statusCode int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteResponse(w http.ResponseWriter, statusCode int, message string, data any) {
	WriteJSON(w, statusCode, APIResponse{
		Status:  statusCode >= 200 && statusCode < 300,
		Message: message,
		Data:    data,
	})
}
