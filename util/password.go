package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("failed to hash password : %w", err)
	}
	return string(hashPassword), nil
}

func CheckPassword(password string, hashedpassword string) error {

	return bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
}
