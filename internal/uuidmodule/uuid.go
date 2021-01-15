package uuidmodule

import (
	"github.com/google/uuid"
)

// Uniquetoken() ...
func Uniquetoken() string {
	unique_id := uuid.New()
	return unique_id.String()
}
