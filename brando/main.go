package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

var (
	dbService, _  = getService()
	defaultTable  = os.Getenv("DB_NAME")
	defaultRegion = os.Getenv("DB_REGION")
	corsOrigin    = os.Getenv("CORS_ORIGIN")
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
	if request.HTTPMethod == http.MethodOptions {
		return methOptions(request), err
	}
	if request.HTTPMethod != http.MethodGet && !authorized(request) {
		return sendMsg("Unauthorized", 401), err
	}
	switch request.HTTPMethod {
	case http.MethodGet:
		return methGet(request), err
	case http.MethodPost:
		return methPost(request), err
	case http.MethodDelete:
		return methDelete(request), err
	case http.MethodPut:
		return methPut(request), err
	default:
		return sendMsg("Unsupported HTTP method!", 400), err
	}
}

func sendJSON(v interface{}) events.APIGatewayProxyResponse {
	log.Info("JSON Marshal...")
	data, err := json.Marshal(v)
	if err != nil {
		log.Error(err.Error())
		return sendMsg("Events coudln't be parsed!", 500)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(data),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": corsOrigin,
		},
	}
}

func authorized(request events.APIGatewayProxyRequest) bool {
	auth := request.Headers["Authorization"]

	if auth == "YWRtaW46ZmlsaXBpbm8tZGl2YW4tbGlzdGluZy10ZW5waW4=" {
		return true
	}

	return false
}

func methGet(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var err error

	providesAuth := authorized(request)

	log.Info("GET request on " + request.Path)

	eventID := request.PathParameters["event"]

	if eventID != "" {
		log.Info("Selected event ID: " + eventID)

		fEvents, err := getFullEventByID(dbService, eventID)
		if err != nil {
			log.Error(err.Error())
			return sendMsg("Selected event couldn't be retrieved!", 500)
		}

		if len(*fEvents) == 0 {
			return sendMsg("Selected event wasn't found", 404)
		}

		if providesAuth {
			return sendJSON((*fEvents)[0])
		}

		code, codeOK := request.QueryStringParameters["code"]

		password := (*fEvents)[0].Password

		if password != "" && (!codeOK || password != code) {
			return sendMsg("Please provide a valid access code!", 401)
		}

		fEvent := (*fEvents)[0]

		uEvent := UEvent{
			fEvent.ID,
			Event{
				Title:       fEvent.Title,
				Presentor:   fEvent.Presentor,
				Description: fEvent.Description,
				StartDate:   fEvent.StartDate,
				EndDate:     fEvent.EndDate,
				Live:        fEvent.Live,
				Featured:    fEvent.Featured,
				Teaser:      fEvent.Teaser,
				Source:      fEvent.Source,
				Hidden:      fEvent.Hidden,
			},
		}

		return sendJSON(uEvent)
	}

	log.Info("Retrieving events...")

	if providesAuth {
		fEvents, err := getAllEvents(dbService)
		if err != nil {
			log.Error(err.Error())
			return sendMsg("Events couldn't be retrieved!", 500)
		}

		return sendJSON(fEvents)
	}
	uEvents, err := getEvents(dbService)
	if err != nil {
		log.Error(err.Error())
		return sendMsg("Events couldn't be retrieved!", 500)
	}

	return sendJSON(uEvents)
}

func methPost(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var err error
	var event ProtectedEvent

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

func methDelete(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
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

func methPut(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var err error
	var event ProtectedEvent

	eventID := request.PathParameters["event"]

	if eventID == "" {
		log.Error("No event ID provided! " + eventID)
		return sendMsg("You must provide an event ID!", 400)
	}

	log.Info("JSON Unmarshal...")

	if err = json.Unmarshal([]byte(request.Body), &event); err != nil {
		log.Error(err.Error())
		return sendMsg("Event couldn't be parsed!", 400)
	}

	log.Info("Updating event...")

	if err = putEvent(dbService, eventID, event); err != nil {
		log.Error(err.Error())
		return sendMsg("Event couldn't be updated!", 500)
	}

	return sendMsg("Success!", 201)
}

func methOptions(request events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  corsOrigin,
			"Access-Control-Allow-Methods": "POST, GET, PUT, DELETE, OPTIONS",
			"Access-Control-Allow-Headers": "Authorization, Content-Type",
		},
	}
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
			"Access-Control-Allow-Origin": corsOrigin,
		},
	}
}

func main() {
	lambda.Start(Handler)
}
