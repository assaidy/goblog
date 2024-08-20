package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// ApiFunc is a custom type that defines a function signature returning an error.
type ApiFunc func(http.ResponseWriter, *http.Request) error

// MakeHandlerFunc wraps an ApiFunc and converts it into an http.HandlerFunc.
// It handles errors by converting them to JSON responses and logging the error.
func MakeHandlerFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// Check if the error is an ApiError
			if apiErr, ok := err.(ApiError); ok {
				WriteJSON(w, apiErr.StatusCode, apiErr)
			} else {
				// If not, return a generic internal server error response
				resp := map[string]string{
					"statusCode": "500",
					"msg":        "internal server error",
				}
				WriteJSON(w, http.StatusInternalServerError, resp)
			}
			// Log the error with additional context
			slog.Error("HTTP API error", "err", err.Error(), "path", r.URL.Path)
		}
	}
}

// WriteJSON sends a JSON response with a given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		// If JSON encoding fails, log the error and return it
		slog.Error("Failed to encode JSON response", "err", err.Error())
		return err
	}

	return nil
}

