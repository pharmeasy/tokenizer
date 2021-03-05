package database



import (

	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/db"

)

type DatabaseInterface interface{
	GetSession(TableName string)
    GetItemsByTokenInBatch(tokenIDs [] string) (map[string]db.TokenData, error)
	GetItemsByToken(tokenIDs []string) (map[string]db.TokenData, error)
	PutItem(item db.TokenData) error 
	UpdateMetadataByToken(tokenID string, metadata map[string]string, updatedAt string) error
	
}

