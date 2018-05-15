package main

import (
	"errors"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	defaultDynamoDBTable = "BrandoTable"
	defaultDynamoDBRegion = "eu-west-1"
)

var (
	dbService,_ = getService()
)

func getService() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(defaultDynamoDBRegion)},
	)

	svc := dynamodb.New(sess)

	return svc, err
}

func scanDB() (*[]ULiveEvent, error) {
	params := &dynamodb.ScanInput{
		TableName: aws.String(defaultDynamoDBTable),
		}

	result, err := dbService.Scan(params)
	if err != nil {
		return nil, err
	} 

	obj := []ULiveEvent{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func getEventByID(s string) (*[]ULiveEvent, error) {
	//Define query parameters, search by given ID
	params := &dynamodb.QueryInput{
		TableName: aws.String(defaultDynamoDBTable),
		KeyConditions: map[string]*dynamodb.Condition{
		 "ID": {
		   ComparisonOperator: aws.String("EQ"),
			AttributeValueList:     []*dynamodb.AttributeValue{
			   {
				S: aws.String(s),
				},
			  },
			},
		   },
		 }
	resp, err := dbService.Query(params)
	if err != nil {
		return nil, err
	} 
	liveEvent := []ULiveEvent{}
	err = dynamodbattribute.UnmarshalListOfMaps(resp.Items,  &liveEvent)
	if err == nil {
		return &liveEvent, nil
	}
	return nil, err
}

func putItem(liveEvent LiveEvent) error {
	//Generate new uuid and append it to liveEvent
	uuid, err := newUUID()
	if err != nil {
		return err
	}
	uLiveEvent := ULiveEvent{uuid, liveEvent}

	av, err := dynamodbattribute.MarshalMap(uLiveEvent)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String(defaultDynamoDBTable),
	}
	
	_, err = dbService.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func delItem(eventID string) error {
	item, err := getEventByID(eventID)
	if len(*item) == 0 {
		return errors.New("Item couldn't be found")
	} 
	if err != nil {
		return err
	}
	
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String((*item)[0].ID),
			},
			"date_begin": {
				S: aws.String((*item)[0].date_begin),
			},
		},
		TableName: aws.String(defaultDynamoDBTable),
	}
	
	_, err = dbService.DeleteItem(input)
	
	if err != nil {
		return err
	}
	
   return nil
}
