package main

import (
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
	switch request.HTTPMethod {
    case "GET":
        return show(request)
    case "POST":
        return create(request)
    default:
        return clientError()
    }
}

func show(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "GET",
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}, nil
}

func create(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	var liveEvent LiveEvent
	
	err = json.Unmarshal([]byte(request.Body), &liveEvent)
	if err != nil {
		err := Error{"Inconsistent input", err.Error()}
		data, _ := json.Marshal(err);
		return events.APIGatewayProxyResponse{Body: string(data), StatusCode: 400}, nil
	}

	av, err := dynamodbattribute.MarshalMap(liveEvent)

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String("BrandoTable"),
	}
	
	_, err = svc.PutItem(input)
	
	if err != nil {
		err := Error{"Could not insert item into database", err.Error()}
		data, _ := json.Marshal(err);
		return events.APIGatewayProxyResponse{Body: string(data), StatusCode: 400}, nil
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func clientError() (events.APIGatewayProxyResponse, error) {
    return events.APIGatewayProxyResponse{
        StatusCode: 404,
        Body:       "ERROR",
    }, nil
}

func main() {
	lambda.Start(Handler)
}
