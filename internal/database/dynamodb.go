package database

import (
	"fmt"

	"bitbucket.org/pharmaeasyteam/goframework/logging"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pharmaeasy/tokenizer/internal/errormanager"
	"github.com/pharmaeasy/tokenizer/internal/models/db"
)

var dbSession *dynamodb.DynamoDB
var tableName string

/*Interface implementation*/
type DynamoDbObject struct {
	TableName string
}

func GetDynamoDbObject(tableName string) *DynamoDbObject {
	newDynamoObject := DynamoDbObject{TableName: tableName}

	return &newDynamoObject
}

// GetSession creates a session if not present
func (d *DynamoDbObject) GetSession(dynamoTableName string) {
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

func (d *DynamoDbObject) GetItemsByTokenInBatch(tokenIDs []string) (map[string]db.TokenData, error) {
	itemsByTokenIDs := make(map[string]db.TokenData)

	tokenLength := len(tokenIDs)
	var filterArray []map[string]*dynamodb.AttributeValue

	for i := 0; i < tokenLength; i++ {
		filterArray = append(filterArray, map[string]*dynamodb.AttributeValue{
			"TokenID": {
				S: aws.String(tokenIDs[i]),
			},
		})
	}
	isConsistentRead := true
	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: {
				Keys:           filterArray,
				ConsistentRead: &isConsistentRead,
			},
		},
	}

	result, err := dbSession.BatchGetItem(input)

	if err != nil {
		return nil, errormanager.SetError("Error encountered while getting DynamoDB item", err)
	}

	dataList := result.Responses[tableName]
	unprocessedTokens := result.UnprocessedKeys[tableName]

	if unprocessedTokens != nil && len(unprocessedTokens.Keys) > 0 {
		unprocessedTokenList := ""
		for i := 0; i < len(unprocessedTokens.Keys); i++ {
			unprocessedTokenList = unprocessedTokenList + *unprocessedTokens.Keys[i]["TokenID"].S
		}
		return nil, errormanager.SetError(fmt.Sprintf("All keys not processed.Token list : %s", unprocessedTokenList), nil)
	}

	resultLen := len(dataList)

	for i := 0; i < resultLen; i++ {

		item := db.TokenData{}

		err = dynamodbattribute.UnmarshalMap(dataList[i], &item)
		if err != nil {
			return nil, errormanager.SetError("Failed to unmarshal record", err)
		}

		itemsByTokenIDs[item.TokenID] = item
	}

	if len(dataList) != len(filterArray) {
		return nil, errormanager.SetError("All DynamoDB Item not found", nil)
	}

	return itemsByTokenIDs, nil

}

// GetItemsByToken Gets the token record from the db
func (d *DynamoDbObject) GetItemsByToken(tokenIDs []string) (map[string]db.TokenData, error) {
	itemsByTokenIDs := make(map[string]db.TokenData)
	isConsistentRead := true
	for _, tokenID := range tokenIDs {
		result, err := dbSession.GetItem(&dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"TokenID": {
					S: aws.String(tokenID),
				},
			},
			ConsistentRead: &isConsistentRead,
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
			return nil, errormanager.SetError(fmt.Sprintf("Failed to unmarshal record, %v for tokenID %s", err, tokenID), err)
		}

		itemsByTokenIDs[tokenID] = item

	}

	return itemsByTokenIDs, nil
}

// PutItem stores the record in the db
func (d *DynamoDbObject) PutItem(item db.TokenData) error {

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return errormanager.SetError(fmt.Sprintf("Failed to unmarshal record for PutItem with tokenID %s", item.TokenID), err)
	}

	input := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(tableName),
		ConditionExpression: aws.String("attribute_not_exists(TokenID)"),
	}

	_, err = dbSession.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

// UpdateMetadataByToken updated attributes in the existing record
func (d *DynamoDbObject) UpdateMetadataByToken(tokenID string, metadata map[string]string, updatedAt string) error {

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
			"TokenID": {
				S: aws.String(tokenID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set Metadata = :metadata, UpdatedAt = :updatedAt"),
	}

	_, err := dbSession.UpdateItem(input)
	if err != nil {
		return errormanager.SetError(fmt.Sprintf("Failed to execute DynamoDB UpdateItem for tokenID %s", tokenID), err)
	}

	return nil
}

// DeleteItemByToken deletes an item from dynamodb

func (d *DynamoDbObject) DeleteItemByToken(tokenID string) error {
	ALL_OLD := "ALL_OLD"
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"TokenID": {
				S: aws.String(tokenID),
			},
		},
		ReturnValues: &ALL_OLD,
	}

	output, err := dbSession.DeleteItem(input)
	deletedTokenId := output.Attributes["TokenID"]

	if *deletedTokenId.S == tokenID {
		logging.GetLogger().Info(fmt.Sprintf("Succesffuly deleted %s", tokenID))
	}
	if err != nil {
		return errormanager.SetError(fmt.Sprintf("Failed to delete the item for tokenID %s", tokenID), err)
	}

	return nil
}
