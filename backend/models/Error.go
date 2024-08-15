package models

import "fmt"

// DBError represents a generic database error.
type DBError struct {
	Err error
}

func (e *DBError) Error() string {
	return fmt.Sprintf("database error: %v", e.Err)
}

func NewDBError(err error) error {
	return &DBError{Err: err}
}

// NotFoundError represents a not found error.
type NotFoundError struct {
	Resource string
	ID       int
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}

func NewNotFoundError(resource string, id int) error {
	return &NotFoundError{Resource: resource, ID: id}
}

// notfoundError represents a bad request error like sending a non-integer id.
type BadRequest struct {
	Err error
}

func (e *BadRequest) Error() string {
	return fmt.Sprintf("bad request: %v", e.Err)
}

func NewBadRequestError(err error) error {
	return &BadRequest{Err: err}
}

// JSONEncodeError represents an error that occurs during JSON encoding.
type JSONEncodeError struct {
	Err error
}

func (e *JSONEncodeError) Error() string {
	return fmt.Sprintf("failed to encode response to JSON: %v", e.Err)
}

func NewJSONEncodeError(err error) error {
	return &JSONEncodeError{Err: err}
}
