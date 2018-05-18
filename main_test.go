package main

// import (
// 	"encoding/json"
// 	"net/http"

// 	"testing"

// 	"github.com/aws/aws-lambda-go/events"
// 	"github.com/stretchr/testify/assert"
// )

// func TestHandler(t *testing.T) {
// 	request := events.APIGatewayProxyRequest{}

// 	data, _ := json.Marshal(Event{
// 		Titel:       "TitelX1",
// 		Presentor:   "PresentorX1",
// 		Description: "DescriptionX1",
// 		StartDate:   "StartDateX1",
// 		EndDate:     "EndDateX1",
// 		Live:        true,
// 		Featured:    false,
// 	})

// 	request.HTTPMethod = http.MethodPost
// 	request.Body = string(data)

// 	responseBody, _ := json.Marshal(Msg{
// 		Message: "Success!",
// 	})

// 	expectedResponse := events.APIGatewayProxyResponse{
// 		StatusCode: 200,
// 		Body:       string(responseBody),
// 		Headers: map[string]string{
// 			"Content-Type": "application/json",
// 		},
// 	}

// 	response := handler(request)

// 	assert.Contains(t, response.Body, expectedResponse)
// }
