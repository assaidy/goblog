package models

import (
	"time"
)

// User represents a user entity in the system.
// This struct is used to store and retrieve user information from the database.
// Fields like `Password` should be handled securely (e.g., hashed and not exposed in responses).
type User struct {
	Id       int       `json:"id"`
	FullName string    `json:"fullName"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Bio      string    `json:"bio"`
	JoinedAt time.Time `json:"joinedAt"`
}

// UserRegisterOrUpdateRequest is used for creating or updating a user account.
// It contains the data required from the client to register or modify a user's information.
type UserRegisterOrUpdateRequest struct {
	FullName string `json:"fullName"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

// UserLoginRequest is used for user login requests.
// It contains the essential fields required to authenticate a user.
type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

