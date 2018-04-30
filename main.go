package main

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

//LiveEvent Struct
type LiveEvent struct {
	Titel string`json:"Titel"`
	Presentor string`json:"Presentor"`
	Description string`json:"Description"`
    	DateBegin string`json:"DateBegin"`
	DateEnd string`json:"DateEnd"`
	Live bool`json:"Live"`
	Featured bool`json:"Featured"`
}

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Println("Received body: ", request.Body)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	var liveEvent LiveEvent
	
	err = json.Unmarshal([]byte(request.Body), &liveEvent)
	if err != nil {
		fmt.Println("Got error parsing JSON:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{Body: "Inconsistent input", StatusCode: 400}, nil
	}
	fmt.Println(liveEvent, err)

	av, err := dynamodbattribute.MarshalMap(liveEvent)

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String("BrandoTable"),
	}
	
	_, err = svc.PutItem(input)
	
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{Body: "Could not insert item into database", StatusCode: 400}, nil
	}
	
	fmt.Println("Successfully added 'The Big New Movie' (2015) to Movies table")

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil

}

func main() {
	lambda.Start(Handler)
}
