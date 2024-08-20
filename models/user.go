package models

import "time"

// User represents a user entity in the system.
// This struct is used to store and retrieve user information from the database.
// Fields like `Password` should be handled securely (e.g., hashed and not exposed in responses).
type User struct {
	Id       int       `json:"id"`       // Unique identifier for the user.
	FullName string    `json:"fullName"` // Full name of the user.
	Username string    `json:"username"` // Unique username for login and display.
	Email    string    `json:"email"`    // User's email address.
	Password string    `json:"password"` // User's password (usually stored as a hashed value).
	Bio      string    `json:"bio"`      // Short biography or description of the user.
	JoinedAt time.Time `json:"joinedAt"` // Timestamp when the user joined the platform.
}

// UserRegisterOrUpdateRequest is used for creating or updating a user account.
// It contains the data required from the client to register or modify a user's information.
type UserRegisterOrUpdateRequest struct {
	FullName string `json:"fullName"` // Full name of the user.
	Username string `json:"username"` // Username for the user account.
	Email    string `json:"email"`    // Email address of the user.
	Password string `json:"password"` // Password for the user account.
	Bio      string `json:"bio"`      // Optional short biography or description of the user.
}

// UserLoginRequest is used for user login requests.
// It contains the essential fields required to authenticate a user.
type UserLoginRequest struct {
	Username string `json:"username"` // Username or email address used for login.
	Password string `json:"password"` // Password for authentication.
}

