package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandler(t *testing.T) {

	request := events.APIGatewayProxyRequest{}
	data, _ := json.Marshal(LiveEvent{Titel: "Awesome Event", Presentor: "Bob", Description: "An awesome event I guess", DateBegin: "2018-05-01 12:00", DateEnd: "2018-05-01 12:30", Live: true, Featured: true })
	request.HTTPMethod = "POST"
	request.Body = string(data)
	expectedResponse := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body: string(data),
	}

	response, err := Handler(request)

	assert.Contains(t, response.Body, expectedResponse.Body)
	assert.Equal(t, err, nil)

}
