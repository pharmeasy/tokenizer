package db

// TokenData represents the struct that stores token related data in DynamoDb
type TokenData struct {
	Token	string 	
	Level 	int
	ARN		string
	Content	string
}

