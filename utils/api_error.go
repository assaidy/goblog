package utils

import (
	"fmt"
	"net/http"
)

// ApiError represents a structured API error with an HTTP status code and a message
type ApiError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"` // Can be a string or a slice of error messages
}

// Error implements the error interface for ApiError
func (e ApiError) Error() string {
	return fmt.Sprintf("api error: %d - %v", e.StatusCode, e.Msg)
}

// NewApiError creates a new ApiError with the given status code and error message
func NewApiError(statusCode int, err error) ApiError {
	return ApiError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

// InvalidRequestData returns an ApiError for invalid request data with a 422 status code
func InvalidRequestData(errors []string) ApiError {
	return ApiError{
		StatusCode: http.StatusUnprocessableEntity, // 422
		Msg:        errors,
	}
}

// InvalidJSON returns an ApiError for invalid JSON request data with a 400 status code
func InvalidJSON() ApiError {
	return NewApiError(http.StatusBadRequest, fmt.Errorf("invalid JSON request data"))
}

// NotFound returns an ApiError for resource not found with a 404 status code
func NotFound(err error) ApiError {
	return NewApiError(http.StatusNotFound, err)
}

// UnAuthorized returns an ApiError for unauthorized access with a 401 status code
func UnAuthorized(err error) ApiError {
	return NewApiError(http.StatusUnauthorized, err)
}

