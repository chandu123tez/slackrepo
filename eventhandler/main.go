package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
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
	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(request.Body),
		slackevents.OptionVerifyToken(
			&slackevents.TokenComparator{VerificationToken: "your-secret-slack-signing-token"},
		),
	)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(request.Body), &r)
		if err != nil {
			return events.APIGatewayProxyResponse{}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       r.Challenge,
		}, nil
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch innerEvent.Data.(type) {
		case *slackevents.AppHomeOpenedEvent:
			event := innerEvent.Data.(*slackevents.AppHomeOpenedEvent)
			// Do things here			
		case *slackevents.MessageEvent:
			m := innerEvent.Data.(*slackevents.MessageEvent)
			// Do things here
		case *slackevents.AppMentionEvent:
			m := innerEvent.Data.(*slackevents.AppMentionEvent)
			// Do things here
		}

	}
	resp := events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "pupster", // Your service name
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
