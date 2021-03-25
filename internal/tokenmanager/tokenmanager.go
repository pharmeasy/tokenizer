package tokenmanager

import (
	"fmt"

	"github.com/google/uuid"
)

var instanceID string

// Uniquetoken generates a unique token for the request
func Uniquetoken() string {
	return uuid.New().String()
}

// FormatToken formats the token response
func FormatToken(token string) string {
	tokenformat := "token://%s/%s"
	formattedToken := fmt.Sprintf(tokenformat, instanceID, token)

	return formattedToken
}

// LoadInstanceIDFromConfig loads instance id from env
func LoadInstanceIDFromConfig(str string) {
	instanceID = str
}
