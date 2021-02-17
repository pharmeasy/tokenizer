package tokenmanager

import (
	"github.com/google/uuid"
)

// Uniquetoken() ...
func Uniquetoken() string {
	unique_id := uuid.New()
	return unique_id.String()
}

func FormatToken(token string) string {
	appendString := "token://a1/"
	appendString = appendString + token
	return appendString
}
