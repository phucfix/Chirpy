package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashBcryptByte, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return "", fmt.Errorf("Can not generate from password: %w", err)
	}

	return string(hashBcryptByte), nil
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
