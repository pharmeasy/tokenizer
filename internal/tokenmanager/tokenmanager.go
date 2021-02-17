package tokenmanager

import (
	"fmt"

	"github.com/google/uuid"
)

// Uniquetoken generates a unique token for the request
func Uniquetoken() string {
	return uuid.New().String()
}

// FormatToken formats the token response
func FormatToken(token string) string {
	tokenformat := "token://%s/%s"
	instanceID := "A1"
	formattedToken := fmt.Sprintf(tokenformat, instanceID, token)

	return formattedToken
}
