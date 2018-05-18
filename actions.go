package main

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	defaultTable  = "BrandoTable"
	defaultRegion = "eu-west-1"
)

var (
	dbService, _ = getService()
)

func getService() (*dynamodb.DynamoDB, error) {
	var err error

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(defaultRegion),
	})

	if err != nil {
		return nil, err
	}

	svc := dynamodb.New(sess)

	return svc, err
}

func getEvents() (*[]UniqueEvent, error) {
	var err error

	params := &dynamodb.ScanInput{
		TableName: aws.String(defaultRegion),
	}

	result, err := dbService.Scan(params)

	if err != nil {
		return nil, err
	}

	events := []UniqueEvent{}

	if err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &events); err != nil {
		return nil, err
	}

	return &events, err
}

func getEventByID(eventID string) (*[]UniqueEvent, error) {
	var err error
	events := []UniqueEvent{}

	params := &dynamodb.QueryInput{
		TableName: aws.String(defaultTable),
		KeyConditions: map[string]*dynamodb.Condition{
			"ID": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(eventID),
					},
				},
			},
		},
	}

	result, err := dbService.Query(params)

	if err != nil {
		return nil, err
	}

	if err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &events); err != nil {
		return nil, err
	}

	return &events, err
}

func putEvent(event Event) error {
	var err error
	var uEvent UniqueEvent
	uuid, err := newUUID()

	if err != nil {
		return err
	}

	uEvent = UniqueEvent{uuid, event}

	item, err := dynamodbattribute.MarshalMap(uEvent)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(defaultTable),
	}

	_, err = dbService.PutItem(input)

	if err != nil {
		return err
	}

	return err
}

func deleteEvent(eventID string) error {
	var err error

	item, err := getEventByID(eventID)

	if len(*item) == 0 {
		err = errors.New("Item couldn't be found")
	}
	if err != nil {
		return err
	}

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String((*item)[0].ID),
			},
			"StartDate": {
				S: aws.String((*item)[0].StartDate),
			},
		},
		TableName: aws.String(defaultTable),
	}

	_, err = dbService.DeleteItem(input)

	if err != nil {
		return err
	}

	return err
}
