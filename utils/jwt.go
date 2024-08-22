package utils

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/assaidy/goblog/models"
	"github.com/assaidy/goblog/repo"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware checks if the request contains a valid JWT token and adds the user ID to the request context.
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Remove the "Bearer " prefix from the token string
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Verify the token
		userId, err := verifyTokenAndGetUserID(tokenString)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
			return
		}

		// Add the user ID to the request context
		ctx := context.WithValue(r.Context(), "userId", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// verifyTokenAndGetUserID verifies the JWT token and extracts the user ID from it.
func verifyTokenAndGetUserID(tokenString string) (int, error) {
	config, err := LoadConfig()
	if err != nil {
		return 0, fmt.Errorf("failed to load config")
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	// Check if there was an error during parsing or if the token is invalid
	if err != nil || !token.Valid {
		return 0, UnAuthorized(fmt.Errorf("invalid token: %v", err))
	}

	// Extract the claims (the payload of the token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, UnAuthorized(fmt.Errorf("invalid token claims"))
	}

	// Extract the user ID from the claims
	userID, ok := claims["userId"].(float64)
	if !ok {
		return 0, UnAuthorized(fmt.Errorf("invalid token claims"))
	}

	return int(userID), nil
}

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
		"userId": id,
		"exp":    expirationTime,
	})

	// Sign the token using the configured JWT secret
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// TODO: apply to the code
// getUserIDFromContext retrieves the user ID from the request context.
func GetUserIDFromContext(r *http.Request) (int, error) {
	userId, ok := r.Context().Value("userId").(int)
	if !ok {
		return 0, UnAuthorized(fmt.Errorf("user ID missing or invalid"))
	}
	return userId, nil
}
