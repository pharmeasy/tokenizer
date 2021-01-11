package models

import (
	"github.com/google/uuid"
)

type DataStore struct {
	TokenID				uuid.UUID	
	EncryptedData		[]byte
	Source 				string
	EncryptionMode		int
	Severity			int
}