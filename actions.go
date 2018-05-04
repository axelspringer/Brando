package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/aws/session"
)

func getSession() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	svc := dynamodb.New(sess)

	return svc, err
}

func scanDB() (*[]ULiveEvent, error) {
	svc, err := getSession()

	params := &dynamodb.ScanInput{
		TableName: aws.String("BrandoTable"),
		}

	result, err := svc.Scan(params)
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
	svc, err := getSession()

	//Define query parameters, search by given ID
	params := &dynamodb.QueryInput{
		TableName: aws.String("BrandoTable"),
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
	resp, err := svc.Query(params)
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
	svc, err := getSession()

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
		TableName: aws.String("BrandoTable"),
	}
	
	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
