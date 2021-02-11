package db

// TokenData represents the struct that stores token related data in DynamoDb
type TokenData struct {
	TokenID   string `json:"tokenId"`
	Level     int
	Content   string
	CreatedAt string
	UpdatedAt string
	Key       string
	Meta      string
}
