// utils/utils.go
package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GetJWTSecret retrieves the JWT secret from environment variables
func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default secret if not set, but it's recommended to set it in .env
		return []byte("your_jwt_secret")
	}
	return []byte(secret)
}

// GenerateToken generates a JWT token for a given user ID
func GenerateToken(userID uint) (string, error) {
	// Set token claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(GetJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken parses a JWT token and returns the user ID
func ParseToken(tokenStr string) (uint, error) {
	// Parse token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return GetJWTSecret(), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	// Extract user_id
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id not found in token")
	}

	return uint(userIDFloat), nil
}

// CreateDirIfNotExists creates a directory if it does not exist
func CreateDirIfNotExists(dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		// Create directory with 0755 permissions
		return os.MkdirAll(dir, 0755)
	}
	return err
}
