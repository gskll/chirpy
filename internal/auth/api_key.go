package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("No Authorization header")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || parts[0] != "ApiKey" {
		return "", fmt.Errorf("Invalid Authorization header")
	}
	return parts[1], nil
}
