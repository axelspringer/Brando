package main

import (
	"fmt"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

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
		err := Error{"Inconsistent input", err.Error()}
		data, _ := json.Marshal(err);
		return events.APIGatewayProxyResponse{Body: string(data), StatusCode: 400}, nil
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
		err := Error{"Could not insert item into database", err.Error()}
		data, _ := json.Marshal(err);
		return events.APIGatewayProxyResponse{Body: string(data), StatusCode: 400}, nil
	}
	
	fmt.Println("Successfully added 'The Big New Movie' (2015) to Movies table")

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil

}

func main() {
	lambda.Start(Handler)
}
