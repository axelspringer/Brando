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
// 		StartDate:   "Fri Jun 01 2018",
// 		EndDate:     "Fri Jun 01 2018",
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

// 	response, err := Handler(request)

// 	assert.Contains(t, response.Body, expectedResponse.Body)
// 	assert.Equal(t, err, nil)

// 	data, _ = json.Marshal(Event{
// 		Titel:       "TitelX55",
// 		Presentor:   "PresentorX2",
// 		Description: "DescriptionX1",
// 		StartDate:   "Fri Jun 01 2018",
// 		EndDate:     "Fri Jun 02 2018",
// 		Live:        true,
// 		Featured:    false,
// 	})

// 	request.HTTPMethod = http.MethodPut
// 	request.Body = string(data)
// 	request.PathParameters = make(map[string]string)
// 	request.PathParameters["event"] = "f436053e-e83a-4c34-808f-b061c9e8a980"

// 	responseBody, _ = json.Marshal(Msg{
// 		Message: "Success!",
// 	})

// 	expectedResponse = events.APIGatewayProxyResponse{
// 		StatusCode: 200,
// 		Body:       string(responseBody),
// 		Headers: map[string]string{
// 			"Content-Type": "application/json",
// 		},
// 	}

// 	response, err = Handler(request)

// 	assert.Contains(t, response.Body, expectedResponse.Body)
// 	assert.Equal(t, err, nil)
// }
