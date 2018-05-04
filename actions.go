package main

import (
	"strconv"

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

func scanDB() (*[]LiveEvent, error) {
	svc, err := getSession()

	params := &dynamodb.ScanInput{
		TableName: aws.String("BrandoTable"),
		}

	result, err := svc.Scan(params)
	if err != nil {
		return nil, err
	} 

	obj := []LiveEvent{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &obj)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func getEventByID(i int) (*[]LiveEvent, error) {
	svc, err := getSession()

	params := &dynamodb.QueryInput{
		TableName: aws.String("BrandoTable"),
		KeyConditions: map[string]*dynamodb.Condition{
		 "ID": {
		   ComparisonOperator: aws.String("EQ"),
			AttributeValueList:     []*dynamodb.AttributeValue{
			   {
				N: aws.String(strconv.Itoa(i)),
				},
			  },
			},
		   },
		 }
	resp, err := svc.Query(params)
	if err != nil {
		return nil, err
	} 
	liveEvent := []LiveEvent{}
	err = dynamodbattribute.UnmarshalListOfMaps(resp.Items,  &liveEvent)
	if err == nil {
		return &liveEvent, nil
	}
	return nil, err
}

func putItem(liveEvent LiveEvent) error {
	svc, err := getSession()

	av, err := dynamodbattribute.MarshalMap(liveEvent)
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
