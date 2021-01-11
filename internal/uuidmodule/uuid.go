package uuidmodule

import (
	"github.com/google/uuid"
)

func Uniquetoken() uuid.UUID {
	unique_id := uuid.New()
	return unique_id
}