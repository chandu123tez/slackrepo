package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	code, _ := request.QueryStringParameters["code"]	
	clientID := "your-slack-client-id"
	clientSecret := "your-slack-client-secret"
	oAuthV2Response, err := slack.GetOAuthV2Response(http.DefaultClient, clientID, clientSecret, code, "")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	client := slack.New(oAuthV2Response.AccessToken)
	teamInfo, err := client.GetTeamInfo()
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	resp := events.APIGatewayProxyResponse{
		StatusCode:      302,
		IsBase64Encoded: false,
		Body:            "",
		Headers: map[string]string{
			"Content-Type":           "text/html",
			"X-MyCompany-Func-Reply": "pupster", // your service name
			"Location":               fmt.Sprintf("https://%v.slack.com", teamInfo.Domain),
		},
	}
	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
