package uuidmodule

import (
	"github.com/google/uuid"
)

// Uniquetoken() ...
func Uniquetoken() string {
	unique_id := uuid.New()
	return append(unique_id.String())
}

func append(token string) string {
	appendString := "token://a1/"
	appendString = appendString + token
	return appendString
}
