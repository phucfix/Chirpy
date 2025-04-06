package auth

import (
	"strings"
	"errors"
	"fmt"
	"net/http"

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

/*
When the user wants to make a request to the API, they send the token along with the request in the HTTP headers. The server can then verify that the token is valid, which means the user is who they say they are.
*/
func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("Authorization header doesn't exist")
	}

	// Split "Bearer" and token
	parts := strings.Fields(auth)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Invalid Authorization header format\n")
	} 

	// The token is the second part
	return parts[1], nil
}
