package models

import "fmt"

// dbError represents a generic database error.
type dbError struct {
	Err error
}

func (e *dbError) Error() string {
	return fmt.Sprintf("database error: %v", e.Err)
}

func NewDBError(err error) error {
	return &dbError{Err: err}
}

// notfoundError represents a not found error.
type notfoundError struct {
	Resource string
	ID       int
}

func (e *notfoundError) Error() string {
	return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}

func NewNotFoundError(resource string, id int) error {
	return &notfoundError{Resource: resource, ID: id}
}

// notfoundError represents a bad request error like sending a non-integer id.
type badRequest struct {
	Err error
}

func (e *badRequest) Error() string {
	return fmt.Sprintf("bad request: %v", e.Err)
}

func NewBadRequestError(err error) error {
	return &badRequest{Err: err}
}
