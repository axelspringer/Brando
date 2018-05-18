package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang/glog"
)

func handler(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	switch request.HTTPMethod {
	case http.MethodGet:
		return get(request)
	case http.MethodPost:
		return post(request)
	case http.MethodDelete:
		return delete(request)
	default:
		return sendMsg("Unsupported HTTP method!", 400)
	}
}

func get(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var err error
	var items *[]UniqueEvent

	glog.Info("GET request on " + request.Path)

	if eventID := request.PathParameters["event"]; eventID != "" {

		glog.Info("Selected event ID: " + eventID)

		items, err = getEventByID(eventID)

		if err != nil {
			glog.Error(err.Error())
			return sendMsg("Selected event couldn't be retrieved!", 500)
		}
	} else {

		glog.Info("Retrieving events...")

		items, err = getEvents()

		if err != nil {
			glog.Error(err.Error())
			return sendMsg("Events couldn't be retrieved!", 500)
		}
	}

	glog.Info("JSON Marshal...")

	data, err := json.Marshal(items)

	if err != nil {
		glog.Error(err.Error())
		return sendMsg("Events coudln't be parsed!", 500)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(data),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func post(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var err error
	var event Event

	glog.Info("JSON Unmarshal...")

	if err = json.Unmarshal([]byte(request.Body), &event); err != nil {
		glog.Error(err.Error())
		return sendMsg("Event couldn't be parsed!", 400)
	}

	if err = putEvent(event); err != nil {
		glog.Error(err.Error())
		return sendMsg("Event couldn't be put into database!", 500)
	}

	return sendMsg("Success!", 200)
}

func delete(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var err error
	eventID := request.PathParameters["event"]

	if eventID == "" {
		glog.Error("No event ID provided! " + eventID)
		return sendMsg("You must provide an event ID!", 400)
	}

	glog.Info("Deleting " + eventID + "...")

	if err = deleteEvent(eventID); err != nil {
		glog.Error(err.Error())
		return sendMsg("Event couldn't be deleted!", 500)
	}

	return sendMsg("Success!", 200)
}

func sendMsg(msg string, status int) events.APIGatewayProxyResponse {
	data, _ := json.Marshal(Msg{
		Message: msg,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(data),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func main() {
	lambda.Start(handler)
}
