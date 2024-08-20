package models

// TODO: apply to the code
// Request is an interface for different request types.
type Request interface {
	Validate() []string
}
