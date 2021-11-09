package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"github.com/slack-go/slack"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	sv, err := slack.NewSecretsVerifier(request.MultiValueHeaders, "your-signing-token")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if _, err := sv.Write([]byte(request.Body)); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if err := return sv.Ensure(); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	options := []*slack.OptionBlockObject{}
	for _, thing := range things {
		option := &slack.OptionBlockObject{
			Text:        slack.NewTextBlockObject(slack.PlainTextType, thing.Text, false, false),
			Description: slack.NewTextBlockObject(slack.PlainTextType, thing.Description, false, false),
			Value:       thing.Value,
		}
		options = append(options, option)		
	}
	optionsResponse := &slack.OptionsResponse{
		Options: options,
	}
	optionsResponseBytes, err := json.Marshal(optionsResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            string(optionsResponseBytes),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "pupster", // your service name
		},
	}, nil
}

func main() {
	lambda.Start(Handler)
}
