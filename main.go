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
	switch request.HTTPMethod {
    case "GET":
        return show(request)
    case "POST":
        return create(request)
    default:
		return clientError(
			Error{"Unsupported http method",
			""}, 404)
    }
}

func getSession() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)
	
	// Create DynamoDB client
	svc := dynamodb.New(sess)

	return svc, err
}

func show(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	svc, err := getSession()

	params := &dynamodb.ScanInput{
		TableName: aws.String("BrandoTable"),
		}

	result, err := svc.Scan(params)
	if err != nil {
	fmt.Printf("failed to make Query API call, %v", err)
	} 

	obj := []LiveEvent{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &obj)
	if err != nil {
	fmt.Printf("failed to unmarshal Query result items, %v", err)
	}

	data, err := json.Marshal(obj)
	if err != nil {
        panic (err)
    }

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(data),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func create(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	
	svc, err := getSession()

	var liveEvent LiveEvent
	
	//unmarshal request body to LiveEvent obj
	err = json.Unmarshal([]byte(request.Body), &liveEvent)
	if err != nil {
		err := Error{"Inconsistent input", err.Error()}
		return clientError(err, 400)
	}

	av, err := dynamodbattribute.MarshalMap(liveEvent)

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String("BrandoTable"),
	}
	
	_, err = svc.PutItem(input)
	
	if err != nil {
		err := Error{"Could not insert item into database", err.Error()}
		return clientError(err, 500)
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func clientError(err Error, code int) (events.APIGatewayProxyResponse, error) {
	data, _ := json.Marshal(err);
    return events.APIGatewayProxyResponse{
        StatusCode: code,
        Body:       string(data),
    }, nil
}

func main() {
	lambda.Start(Handler)
}
