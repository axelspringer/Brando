package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandler(t *testing.T) {

	request := events.APIGatewayProxyRequest{}
	data, _ := json.Marshal(LiveEvent{titel: "Awesome EventXXX", presentor: "Bob", description: "An awesome event I guess", date_begin: "2018-05-01 12:00", date_end: "2018-05-01 12:30", live: true, featured: true })
	request.HTTPMethod = http.MethodPost
	request.Body = string(data)
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: string(data),
	}

	response, err := Handler(request)

	assert.Contains(t, response.Body, expectedResponse.Body)
	assert.Equal(t, err, nil)

	// request = events.APIGatewayProxyRequest{}
	// data, _ = json.Marshal(LiveEventID{ID: "52691d2f-e2da-426f-ac64-6d503edf90ea" })
	// request.HTTPMethod = http.MethodDelete
	// request.Body = string(data)
	// expectedResponse = events.APIGatewayProxyResponse{
	// 	StatusCode: 200,
	// 	Body: string(data),
	// }

	// response, err = Handler(request)

	// assert.Contains(t, response.Body, expectedResponse.Body)
	// assert.Equal(t, err, nil)

	// request = events.APIGatewayProxyRequest{}
	// request.HTTPMethod = "GET"
	// request.Path = "/events/1"
	// request.Body = ""
	// expectedResponse := events.APIGatewayProxyResponse{
	// 	StatusCode: 200,
	// 	Body: "GET",
	// }

	// response, err := Handler(request)
	// fmt.Println(response.Body)
	// assert.Contains(t, response.Body, expectedResponse.Body)
	// assert.Equal(t, err, nil)

}
