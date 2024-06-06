package database

import (
	"github.com/pharmaeasy/tokenizer/internal/models/db"
)

type DatabaseInterface interface {
	GetSession(TableName string)
	GetItemsByTokenInBatch(tokenIDs []string) (map[string]db.TokenData, error)
	GetItemsByToken(tokenIDs []string) (map[string]db.TokenData, error)
	PutItem(item db.TokenData) error
	UpdateMetadataByToken(tokenID string, metadata map[string]string, updatedAt string) error
	DeleteItemByToken(tokenID string) error
}
