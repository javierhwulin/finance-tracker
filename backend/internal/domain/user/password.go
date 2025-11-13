package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

// HashPassword hashes a plaintext password using bcrypt
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password is required")
	}
	if len(password) < 8 {
		return "", errors.New("password must be at least 8 characters long")
	}
	if len(password) > 72 {
		return "", errors.New("password must be less than 72 characters long")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword compares a plaintext password with a bcrypt hash
func ComparePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
