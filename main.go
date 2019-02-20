package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

const (
	defaultTable  = "BrandoTable"
	defaultRegion = "eu-west-1"
)

var (
	dbService, _ = getService()
)

func getService() (*dynamodb.DynamoDB, error) {
	var err error

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(defaultRegion),
	})

	if err != nil {
		return nil, err
	}

	svc := dynamodb.New(sess)

	return svc, err
}

func init() {
	// Seed the default rand Source with current time to produce better random
	// numbers used with splay
	rand.Seed(time.Now().UnixNano())
}

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var err error
	if request.HTTPMethod != http.MethodGet && !authorized(request) {
		return sendMsg("Unauthorized", 401), err
	}
	switch request.HTTPMethod {
	case http.MethodGet:
		return get(request), err
	case http.MethodPost:
		return post(request), err
	case http.MethodDelete:
		return delete(request), err
	case http.MethodPut:
		return put(request), err
	default:
		return sendMsg("Unsupported HTTP method!", 400), err
	}
}

func authorized(request events.APIGatewayProxyRequest) bool {
	auth := request.Headers["Authorization"]

	if auth == "YWRtaW46ZmlsaXBpbm8tZGl2YW4tbGlzdGluZy10ZW5waW4=" {
		return true
	}

	return false
}

func get(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var err error
	var items *[]UniqueEvent
	var data []byte

	log.Info("GET request on " + request.Path)

	eventID := request.PathParameters["event"]

	if eventID != "" {

		log.Info("Selected event ID: " + eventID)

		items, err = getEventByID(dbService, eventID)

		if err != nil {
			log.Error(err.Error())
			return sendMsg("Selected event couldn't be retrieved!", 500)
		}

		data, err = json.Marshal((*items)[0])
	} else {

		log.Info("Retrieving events...")

		if authorized(request) {
			items, err = getAllEvents(dbService)
		} else {
			items, err = getEvents(dbService)
		}

		if err != nil {
			log.Error(err.Error())
			return sendMsg("Events couldn't be retrieved!", 500)
		}

		data, err = json.Marshal(items)
	}

	log.Info("JSON Marshal...")
	if err != nil {
		log.Error(err.Error())
		return sendMsg("Events coudln't be parsed!", 500)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(data),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}

func post(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var err error
	var event Event

	log.Info("JSON Unmarshal...")

	if err = json.Unmarshal([]byte(request.Body), &event); err != nil {
		log.Error(err.Error())
		return sendMsg("Event couldn't be parsed!", 400)
	}

	log.Info("Putting event into database...")

	if err = postEvent(dbService, event); err != nil {
		log.Error(err.Error())
		log.Error(event)
		return sendMsg("Event couldn't be put into database!", 500)
	}

	return sendMsg("Success!", 201)
}

func delete(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var err error
	eventID := request.PathParameters["event"]

	if eventID == "" {
		log.Error("No event ID provided! " + eventID)
		return sendMsg("You must provide an event ID!", 400)
	}

	log.Info("Deleting " + eventID + "...")

	if err = deleteEvent(dbService, eventID); err != nil {
		log.Error(err.Error())
		return sendMsg("Event couldn't be deleted!", 500)
	}

	return sendMsg("Success!", 200)
}

func put(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var err error
	var event Event

	eventID := request.PathParameters["event"]

	if eventID == "" {
		log.Error("No event ID provided! " + eventID)
		return sendMsg("You must provide an event ID!", 400)
	}

	items, err := getEventByID(dbService, eventID)

	if err != nil || len(*items) == 0 {
		log.Error(err.Error())
		return sendMsg("Event doesn't exist!", 500)
	}

	log.Info("JSON Unmarshal...")

	if err = json.Unmarshal([]byte(request.Body), &event); err != nil {
		log.Error(err.Error())
		return sendMsg("Event couldn't be parsed!", 400)
	}

	log.Info("Deleting old " + eventID + "...")

	if err = deleteEvent(dbService, eventID); err != nil {
		log.Error(err.Error())
		return sendMsg("Event couldn't be deleted!", 500)
	}

	log.Info("Updating event...")

	if err = putEvent(dbService, eventID, event); err != nil {
		log.Error(err.Error())
		log.Error(event)
		return sendMsg("Event couldn't be updated!", 500)
	}

	return sendMsg("Success!", 201)
}

func sendMsg(msg string, status int) events.APIGatewayProxyResponse {
	data, _ := json.Marshal(Msg{
		Message: msg,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       string(data),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
}

func main() {
	lambda.Start(Handler)
}
