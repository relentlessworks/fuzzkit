package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// wantsJSON checks if the client wants JSON response.
func wantsJSON(r *http.Request) bool {
	if r.URL.Query().Get("format") == "json" {
		return true
	}
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "application/json")
}

// writeText writes a plain text response.
func writeText(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprint(w, body)
}

// writeJSON writes a JSON response.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeResponse writes either plain text or JSON based on the request.
func writeResponse(w http.ResponseWriter, r *http.Request, status int, text string, data interface{}) {
	if wantsJSON(r) {
		writeJSON(w, status, data)
	} else {
		writeText(w, status, text)
	}
}

// writeError writes an error response with a hint.
func writeError(w http.ResponseWriter, r *http.Request, status int, message, hint string) {
	if wantsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]string{
			"error": message,
			"hint":  hint,
		})
	} else {
		writeText(w, status, fmt.Sprintf("error: %s | hint: %s\n", message, hint))
	}
}

// writeRecord writes a single key=value record in plain text or JSON.
func writeRecord(w http.ResponseWriter, r *http.Request, fields map[string]interface{}) {
	if wantsJSON(r) {
		writeJSON(w, http.StatusOK, fields)
	} else {
		var parts []string
		for k, v := range fields {
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
		}
		writeText(w, http.StatusOK, strings.Join(parts, " ")+"\n")
	}
}
