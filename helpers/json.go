package helpers

import (
	"encoding/json"
	"net/http"
)

// readJSON decodes the request body into the provided struct
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) {
	maxBytes := 1048576 // 1MB max payload
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {
		http.Error(w, "Error reading JSON", http.StatusBadRequest)
		return
	}
}

// writeJSON encodes the provided data as JSON and writes it to the response
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(data)
	if err != nil {
		http.Error(w, "Error writing JSON", http.StatusInternalServerError)
		return
	}
}