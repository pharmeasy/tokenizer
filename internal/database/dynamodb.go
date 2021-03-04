package database

import (
	"fmt"

	"bitbucket.org/pharmaeasyteam/tokenizer/internal/errormanager"

	"bitbucket.org/pharmaeasyteam/tokenizer/internal/models/db"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

)

var dbSession *dynamodb.DynamoDB
var tableName string

// GetSession creates a session if not present
func GetSession(dynamoTableName string) {
	if dbSession == nil {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		dbSession = dynamodb.New(sess)
	}

	if tableName == "" {
		tableName = dynamoTableName
	}
}

// GetItemsByTokenInBatch

func GetItemsByTokenInBatch(tokenIDs [] string) (map[string]db.TokenData, error){
	itemsByTokenIDs := make(map[string]db.TokenData)

	tokenLength := len(tokenIDs)
	var filterArray[] map[string]*dynamodb.AttributeValue

	for i:=0;i<tokenLength;i++{
		filterArray = append(filterArray, map[string]*dynamodb.AttributeValue{
			"tokenId": {
			   S: aws.String(tokenIDs[i]),
			},
		
		})
	}
    
    input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName : {
				Keys: filterArray,
			},
		},
	}

	result, err := dbSession.BatchGetItem(input)

	if err != nil {
		return nil, errormanager.SetError(fmt.Sprintf("Error encountered while getting DynamoDB item"), err)	
	}

	dataList := result.Responses[tableName]
	resultLen := len(dataList)

	for i:=0;i<resultLen;i++{

		item := db.TokenData{}

		err = dynamodbattribute.UnmarshalMap(dataList[i], &item)
		if err != nil {
			return nil, errormanager.SetError(fmt.Sprintf("Failed to unmarshal Record"), err)
		}

		itemsByTokenIDs[item.TokenID] = item
	}

	if len(dataList) != len(filterArray) {
		return nil, errormanager.SetError(fmt.Sprintf("All DynamoDB Item not found"), nil)
	}



    return itemsByTokenIDs, nil



}



// GetItemsByToken Gets the token record from the db
func GetItemsByToken(tokenIDs []string) (map[string]db.TokenData, error) {
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

		if err != nil {
			return nil, errormanager.SetError(fmt.Sprintf("Error encountered while getting DynamoDB item for tokenID %s", tokenID), err)
		}

		// throw 5xx
		if result.Item == nil {
			return nil, errormanager.SetError(fmt.Sprintf("DynamoDB Item not found for tokenID %s", tokenID), nil)
		}

		item := db.TokenData{}

		err = dynamodbattribute.UnmarshalMap(result.Item, &item)
		if err != nil {
			return nil, errormanager.SetError(fmt.Sprintf("Failed to unmarshal Record, %v for tokenID %s", err, tokenID), err)
		}

		itemsByTokenIDs[tokenID] = item

	}

	return itemsByTokenIDs, nil
}

// PutItem stores the record in the db
func PutItem(item db.TokenData) error {

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return errormanager.SetError(fmt.Sprintf("Failed to unmarshal Record for PutItem with tokenID %s", item.TokenID), err)
	}

	input := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_not_exists(tokenId)"),
	}

	_, err = dbSession.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

// UpdateMetadataByToken updated attributes in the existing record
func UpdateMetadataByToken(tokenID string, metadata map[string]string, updatedAt string) error {

	meta, _ := dynamodbattribute.MarshalMap(metadata)
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":metadata": {
				M: meta,
			},
			":updatedAt": {
				S: aws.String(updatedAt),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"tokenId": {
				S: aws.String(tokenID),
			},
		},
		ReturnValues:        aws.String("UPDATED_NEW"),
		UpdateExpression:    aws.String("set Metadata1 = :metadata, UpdatedAt = :updatedAt"),
		ConditionExpression: aws.String("attribute_not_exists(tokenID)"),
	}

	_, err := dbSession.UpdateItem(input)
	if err != nil {
		return errormanager.SetError(fmt.Sprintf("Failed to execute DynamoDB UpdateItem for tokenID %s", tokenID), err)
	}
	//	}

	return nil
}
