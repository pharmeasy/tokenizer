package db

// TokenData represents the struct that stores token related data in DynamoDb
type TokenData struct {
	TokenID   string `json:"tokenId"`
	Level     string
	Content   []byte
	CreatedAt string
	UpdatedAt string
	Key       string
	Metadata1 map[string]string
}
