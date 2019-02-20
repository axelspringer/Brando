package main

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func getEvents(svc dynamodbiface.DynamoDBAPI) (*[]UniqueEvent, error) {
	var err error
	events := &[]UniqueEvent{}
	publicEvents := []UniqueEvent{}

	params := &dynamodb.ScanInput{
		TableName: aws.String(defaultTable),
	}

	result, err := svc.Scan(params)
	if err != nil {
		return events, err
	}

	if err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &events); err != nil {
		return events, err
	}

	for _, event := range *events {
		if !event.Hidden {
			publicEvents = append(publicEvents, event)
		}
	}

	return &publicEvents, err
}

func getAllEvents(svc dynamodbiface.DynamoDBAPI) (*[]UniqueEvent, error) {
	var err error
	events := &[]UniqueEvent{}

	params := &dynamodb.ScanInput{
		TableName: aws.String(defaultTable),
	}

	result, err := svc.Scan(params)
	if err != nil {
		return events, err
	}

	if err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &events); err != nil {
		return events, err
	}

	return events, err
}

func getEventByID(svc dynamodbiface.DynamoDBAPI, eventID string) (*[]UniqueEvent, error) {
	var err error
	events := &[]UniqueEvent{}

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

	result, err := svc.Query(params)
	if err != nil {
		return events, err
	}

	if err = dynamodbattribute.UnmarshalListOfMaps(result.Items, events); err != nil {
		return events, err
	}

	return events, err
}

func postEvent(svc dynamodbiface.DynamoDBAPI, event Event) error {
	var err error
	var uniqueEvent UniqueEvent

	uuid, err := newUUID()
	if err != nil {
		return err
	}

	uniqueEvent = UniqueEvent{uuid, event}

	item, err := dynamodbattribute.MarshalMap(uniqueEvent)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(defaultTable),
		Item:      item,
	}

	if _, err := svc.PutItem(input); err != nil {
		return err
	}

	return err
}

func deleteEvent(svc dynamodbiface.DynamoDBAPI, eventID string) error {
	var err error

	event, err := getEventByID(svc, eventID)
	if err != nil {
		return err
	}

	if len(*event) == 0 {
		return errors.New("Event could not be found")
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(defaultTable),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String((*event)[0].ID),
			},
			"StartDate": {
				S: aws.String((*event)[0].StartDate),
			},
		},
	}

	if _, err = svc.DeleteItem(input); err != nil {
		return err
	}

	return err
}

func putEvent(svc dynamodbiface.DynamoDBAPI, eventID string, event Event) error {
	var err error
	var uniqueEvent UniqueEvent

	oldEvent, err := getEventByID(svc, eventID)
	if err != nil {
		return err
	}

	if len(*oldEvent) == 0 {
		return errors.New("Evebt could not be found")
	}

	uniqueEvent = UniqueEvent{eventID, event}

	item, err := dynamodbattribute.MarshalMap(uniqueEvent)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(defaultTable),
		Item:      item,
	}

	if _, err := svc.PutItem(input); err != nil {
		return err
	}

	return err
}
