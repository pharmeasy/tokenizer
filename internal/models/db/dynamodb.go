package db

// TokenData represents the struct that stores token related data in DynamoDb
type TokenData struct {
	TokenID   string
	Level     string
	Content   []byte
	CreatedAt string
	UpdatedAt string
	Key       string
	Metadata  map[string]string
}
