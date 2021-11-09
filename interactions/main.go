package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var payload slack.InteractionCallback
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
	err := json.Unmarshal([]byte(parseBody(request.Body)), &payload)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	switch payload.Type {
	case slack.InteractionTypeBlockActions:
		// Do things here
	case slack.InteractionTypeViewSubmission:
		// Do things here
	case slack.InteractionTypeShortcut:
		// Do things here
	default:
		err = fmt.Errorf("Unrecognized payload type: %s", payload)
	}
	if err != nil {
		errorText := fmt.Sprintf("Error handling interaction: %s", err.Error())
	}
	return event, err
}

func parseBody(body string) string {
	decodedValue, _ := url.QueryUnescape(body)
	data := strings.Trim(decodedValue, ":payload=")
	return data
}

func main() {
	lambda.Start(Handler)
}
