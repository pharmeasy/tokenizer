package database

import (
	"errors"
	"fmt"

	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var dbSession *dynamodb.DynamoDB

//var tableName string

// GetSession creates a session if not present
func getSession() {
	if dbSession == nil {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		dbSession = dynamodb.New(sess)
	}

	// if tableName == "" {

	// }
}

// GetTableName gets the table name
// func getTableName() {
// 	tableName = "staging_tokens"
// }

// GetItemsByToken Gets the token record from the db
func GetItemsByToken(tokenIDs []string, tableName string) (map[string]db.TokenData, error) {
	getSession()
	itemsByTokenIDs := make(map[string]db.TokenData)

	for _, tokenID := range tokenIDs {
		result, err := dbSession.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"tokenId": {
					S: aws.String(tokenID),
				},
			},
		})

		// throw 5xx
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		// throw 5xx
		if result.Item == nil {
			return nil, errors.New("your request is malformed")
		}

		item := db.TokenData{}

		err = dynamodbattribute.UnmarshalMap(result.Item, &item)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}

		itemsByTokenIDs[tokenID] = item

	}

	return itemsByTokenIDs, nil
}

// PutItem stores the record in the db
func PutItem(item db.TokenData, tableName string) error {
	getSession()

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println(err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_not_exists(tokenId)"),
	}

	_, err = dbSession.PutItem(input)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// UpdateMetadataByToken updated attributes in the existing record
func UpdateMetadataByToken(tokenID string, metadata string, tableName string) error {
	getSession()

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":m": {
				S: aws.String(metadata),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"tokenId": {
				S: aws.String(tokenID),
			},
		},
		ReturnValues:        aws.String("UPDATED_NEW"),
		UpdateExpression:    aws.String("set Metadata = :m"),
		ConditionExpression: aws.String("attribute_not_exists(tokenID)"),
	}

	_, err := dbSession.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}
