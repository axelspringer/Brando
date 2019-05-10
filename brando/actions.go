package main

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func getEvents(svc dynamodbiface.DynamoDBAPI) (*[]UEvent, error) {
	var err error
	events := &[]UEvent{}
	publicEvents := []UEvent{}

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

func getAllEvents(svc dynamodbiface.DynamoDBAPI) (*[]FullEvent, error) {
	var err error
	events := &[]FullEvent{}

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

func getEventByID(svc dynamodbiface.DynamoDBAPI, eventID string) (*[]UEvent, error) {
	var err error
	events := &[]UEvent{}

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

func getFullEventByID(svc dynamodbiface.DynamoDBAPI, eventID string) (*[]FullEvent, error) {
	var err error
	events := &[]FullEvent{}

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

func postEvent(svc dynamodbiface.DynamoDBAPI, event ProtectedEvent) error {
	var err error
	var fEvent FullEvent

	uuid, err := newUUID()
	if err != nil {
		return err
	}

	mEvent := Event{
		Title:       event.Title,
		Presentor:   event.Presentor,
		Description: event.Description,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Live:        event.Live,
		Featured:    event.Featured,
		Teaser:      event.Teaser,
		Source:      event.Source,
		Hidden:      event.Hidden,
	}

	fEvent = FullEvent{uuid, mEvent, event.Password}

	item, err := dynamodbattribute.MarshalMap(fEvent)
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

func putEvent(svc dynamodbiface.DynamoDBAPI, eventID string, event ProtectedEvent) error {
	var err error
	var fEvent FullEvent

	oldEvent, err := getEventByID(svc, eventID)
	if err != nil {
		return err
	}

	if len(*oldEvent) == 0 {
		return errors.New("Event could not be found")
	}

	mEvent := Event{
		Title:       event.Title,
		Presentor:   event.Presentor,
		Description: event.Description,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		Live:        event.Live,
		Featured:    event.Featured,
		Teaser:      event.Teaser,
		Source:      event.Source,
		Hidden:      event.Hidden,
	}

	fEvent = FullEvent{eventID, mEvent, event.Password}

	item, err := dynamodbattribute.MarshalMap(fEvent)
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
