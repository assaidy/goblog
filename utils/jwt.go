package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/assaidy/goblog/models"
	"github.com/assaidy/goblog/repo"
	"github.com/golang-jwt/jwt/v5"
)

// AuthenticateUser returns user data along with a JWT token
func AuthenticateUser(loginReq models.UserLoginRequest, s repo.Storer) (map[string]any, error) {
	user, err := s.GetUserByUsername(loginReq.Username)
	if err != nil {
		return nil, err
	}

	// Check if the password matches
	if user.Password != loginReq.Password {
		return nil, NotFound(fmt.Errorf("password is not correct"))
	}

	// Create a JWT token for the authenticated user
	token, err := createToken(user.Id)
	if err != nil {
		return nil, err
	}

	// Remove the password before returning the user data
	user.Password = ""
	resp := map[string]any{
		"user":  user,
		"token": token,
	}

	return resp, nil
}

// AuthorizeUser checks if the JWT token provided in the request is valid
func AuthorizeUser(id int, r *http.Request) error {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return UnAuthorized(fmt.Errorf("missing authorization header"))
	}

	// Remove "Bearer " prefix from the token string
	tokenString = tokenString[len("Bearer "):]

	// Verify the token
	if err := verifyToken(tokenString); err != nil {
		return err
	}

	return nil
}

// createToken generates a JWT token for a given user ID
func createToken(id int) (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config")
	}

	// Set expiration time based on the configured JWT expiration hours
	expirationTime := time.Now().Add(time.Hour * time.Duration(config.JWTExpirationHours)).Unix()

	// Create a new JWT token with the user ID and expiration time
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": expirationTime,
	})

	// Sign the token using the configured JWT secret
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// verifyToken checks if the provided JWT token string is valid
func verifyToken(tokenString string) error {
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config")
	}

	// Parse the token using the configured JWT secret
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	// Return an unauthorized error if the token is invalid or parsing fails
	if err != nil {
		return UnAuthorized(fmt.Errorf("invalid token: %v", err))
	}

	// Check if the token is valid
	if !token.Valid {
		return UnAuthorized(fmt.Errorf("invalid token"))
	}

	return nil
}
