package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
    case http.MethodGet:
        return show(request)
    case http.MethodPost:
		return create(request)
	case http.MethodDelete:
		return delete(request)
    default:
		return clientError("Unsupported http method", 400)
    }
}


func show(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var data []byte
	var err error

	fmt.Println("GET request on " + request.Path + " event: " + request.PathParameters["event"] + ".")
	if eventID := request.PathParameters["event"]; &eventID != nil && eventID != "" {
		item, err := getEventByID(eventID)
		if err != nil {
			fmt.Println(err.Error())
			return clientError("An unexpected error occured during query", 500)
		}
		data, err = json.Marshal(item)
	} else {
		obj, err := scanDB()
		if err != nil {
			fmt.Println(err.Error())
			return clientError("An unexpected error occured during scan", 500)
		}
		data, err = json.Marshal(obj)
	}

	if err != nil {
		fmt.Println(err.Error())
		return clientError("An unexpected error occured while parsing", 500)
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
		fmt.Println(err.Error())
		return clientError("Inconsistent input", 400)
	}

	if err = putItem(liveEvent); err != nil {
		fmt.Println(err.Error())
		return clientError("An unexpected error occured during put request", 400)
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func delete(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var liveEventID LiveEventID
	
	err := json.Unmarshal([]byte(request.Body), &liveEventID)
	if err != nil {
		fmt.Println(err.Error())
		return clientError("Inconsistent input", 400)
	}

	if err = delItem(liveEventID.ID); err != nil {
		fmt.Println(err.Error())
		return clientError("An unexpected error occured during put request", 400)
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
}

func clientError(errStr string, code int) (events.APIGatewayProxyResponse, error) {
	err := errors.New(errStr)
    return events.APIGatewayProxyResponse{
        StatusCode: code,
        Body:       err.Error(),
    }, nil
}

func main() {
	lambda.Start(Handler)
}
