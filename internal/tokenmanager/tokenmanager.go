package tokenmanager

import (
	"fmt"
	"strings"

	"bitbucket.org/pharmaeasyteam/tokenizer/internal/errormanager"
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

func ExtractToken(token *string) error {
	tokenPrefix := "token://%s/"
	formattedTokenPrefix := fmt.Sprintf(tokenPrefix, instanceID)

	if strings.HasPrefix(*token, formattedTokenPrefix) {
		*token = strings.TrimLeft(*token, formattedTokenPrefix)
		return nil
	} else {
		return errormanager.SetError("Invalid Token Structure", nil)
	}
}
