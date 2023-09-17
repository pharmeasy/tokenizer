package tokenmanager

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pharmaeasy/tokenizer/internal/errormanager"
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
		*token = strings.TrimPrefix(*token, formattedTokenPrefix)
		fmt.Println(*token)
		return nil
	} else {
		return errormanager.SetError("Invalid Token Structure", nil)
	}
}
