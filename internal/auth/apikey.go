package auth

import (
	"net/http"
	"errors"
	"strings"
)

// Extract api key from Authorization header
func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("Authorization header doesn't exist")
	}

	// Split by spaces
	parts := strings.Fields(auth)
	if len(parts) != 2 || parts[0] != "ApiKey" {
		return "", errors.New("Invalid Authorization header format")
	}

	// The api key is in second part
	return parts[1], nil
}
