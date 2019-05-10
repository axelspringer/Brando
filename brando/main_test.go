package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
}

var (
	mockSvc   = &mockDynamoDBClient{}
	mockTable = map[string]FullEvent{}
)

func (m *mockDynamoDBClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	event := FullEvent{}

	dynamodbattribute.UnmarshalMap(input.Item, &event)

	mockTable[event.ID] = event

	return nil, nil
}

func (m *mockDynamoDBClient) Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
	var output []map[string]*dynamodb.AttributeValue

	for _, item := range mockTable {
		marshaled, _ := dynamodbattribute.MarshalMap(item)
		output = append(output, marshaled)
	}

	return &dynamodb.ScanOutput{
		Items: output,
	}, nil
}

func (m *mockDynamoDBClient) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	var output []map[string]*dynamodb.AttributeValue
	id := *(input.KeyConditions["ID"].AttributeValueList)[0].S

	if val, ok := mockTable[id]; ok {
		marshaled, _ := dynamodbattribute.MarshalMap(val)
		output = append(output, marshaled)
		return &dynamodb.QueryOutput{
			Items: output,
		}, nil
	}

	return nil, nil
}

func (m *mockDynamoDBClient) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	id := *input.Key["ID"].S

	delete(mockTable, id)

	return nil, nil
}

func TestMain(t *testing.T) {
	testEvent1 := ProtectedEvent{
		Event{
			Title:       "TestEvent1",
			Presentor:   "Presentor1",
			Description: "Desc1",
			StartDate:   "2019-02-19T09:40:09.699Z",
			EndDate:     "2019-02-19T09:40:09.699Z",
			Live:        true,
			Featured:    false,
			Teaser:      "aW1hZ2U=",
			Source:      "http://example.com",
			Hidden:      true,
		},
		"",
	}

	testEvent2 := ProtectedEvent{
		Event{
			Title:       "TestEvent2",
			Presentor:   "Presentor2",
			Description: "Desc2",
			StartDate:   "2019-02-19T09:40:09.699Z",
			EndDate:     "2019-02-19T09:40:09.699Z",
			Live:        false,
			Featured:    true,
			Teaser:      "aW1hZ2U=",
			Source:      "http://another.com",
			Hidden:      false,
		},
		"",
	}

	testEvent3 := ProtectedEvent{
		Event{
			Title:       "TestEvent3",
			Presentor:   "Presentor3",
			Description: "Desc3",
			StartDate:   "2019-02-19T09:40:09.699Z",
			EndDate:     "2019-02-19T09:40:09.699Z",
			Live:        true,
			Featured:    true,
			Teaser:      "aW1hZ2U=",
			Source:      "http://some.com",
			Hidden:      true,
		},
		"S0mePassw0rd",
	}

	err := postEvent(mockSvc, testEvent1)
	if err != nil {
		t.Fatalf("postEvent (1) errored: %v", err)
	}

	err = postEvent(mockSvc, testEvent2)
	if err != nil {
		t.Fatalf("postEvent (2) errored: %v", err)
	}

	events, err := getEvents(mockSvc)
	if err != nil {
		t.Fatalf("getEvents errored: %v", err)
	}

	for _, event := range *events {
		if event.Hidden {
			t.Fatal("there should be no hidden events returned when calling getEvents")
		}
	}

	id := (*events)[0].ID

	err = deleteEvent(mockSvc, id)
	if err != nil {
		t.Fatalf("deleteEvent errored: %v", err)
	}

	fEvents, err := getAllEvents(mockSvc)
	if err != nil {
		t.Fatalf("GetAllEvents errored: %v", err)
	}

	id = (*fEvents)[0].ID

	err = putEvent(mockSvc, id, testEvent3)
	if err != nil {
		t.Fatalf("putEvent errored: %v", err)
	}

	fEvents, err = getFullEventByID(mockSvc, id)
	if err != nil {
		t.Fatalf("putEvent errored: %v", err)
	}

	if !(*fEvents)[0].Featured || (*fEvents)[0].Password != "S0mePassw0rd" {
		t.Fatalf("seems the update failed!")
	}

}
