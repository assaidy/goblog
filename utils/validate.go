package utils

import (
	"fmt"
	"net/mail"

	"github.com/assaidy/goblog/repo"
)

// ValidateRegisterUser checks if the email is valid and if the username or email are already used.
func ValidateRegisterUser(username, email string, s repo.Storer) (validationErrors []string, internalError error) {
	// Validate email format
	if err := validateEmail(email); err != nil {
		return []string{err.Error()}, nil
	}

	// Check if username or email are already used
	if err := checkUsernameAndEmail(username, email, s, &validationErrors); err != nil {
		return nil, err
	}

	// Return validation errors if any
	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	return nil, nil
}

// validateEmail checks if the email address has a valid format.
func validateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email format: %v", err)
	}
	return nil
}

// checkUsernameAndEmail checks if the username or email are already in use.
func checkUsernameAndEmail(username, email string, s repo.Storer, validationErrors *[]string) error {
	// Check if username is already taken
	if exists, err := s.IsUsernameUsed(username); err != nil {
		return err
	} else if exists {
		*validationErrors = append(*validationErrors, "username is already taken")
	}

	// Check if email is already taken
	if exists, err := s.IsEmailUsed(email); err != nil {
		return err
	} else if exists {
		*validationErrors = append(*validationErrors, "email is already taken")
	}

	return nil
}

