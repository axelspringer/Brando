package main

import (
	"fmt"
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

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


func show(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var data []byte
	var err error

	fmt.Println("GET request on " + request.Path + " event: " + request.PathParameters["event"] + ".")
	if eventID, err := strconv.Atoi(request.PathParameters["event"]); err == nil {
		item, err := getEventByID(eventID)
		if err != nil {
			err := Error{"An unexpected error occured during query", err.Error()}
			return clientError(err, 500)
		}
		data, err = json.Marshal(item)
	} else {
		obj, err := scanDB()
		if err != nil {
			err := Error{"An unexpected error occured during scan", err.Error()}
			return clientError(err, 500)
		}
		data, err = json.Marshal(obj)
	}

	if err != nil {
        err := Error{"An unexpected error occured while parsing", err.Error()}
		return clientError(err, 500)
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
	var liveEvent LiveEvent
	
	// Unmarshal request body to LiveEvent obj
	err := json.Unmarshal([]byte(request.Body), &liveEvent)
	if err != nil {
		err := Error{"Inconsistent input", err.Error()}
		return clientError(err, 400)
	}

	if err = putItem(liveEvent); err != nil {
		err := Error{"An unexpected error occured during put request", err.Error()}
		return clientError(err, 400)
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
