package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type PlausibleInput struct {
	Name        string            `json:"name"`   // e.g. "pageview"
	URL         string            `json:"url"`    // full URL
	Domain      string            `json:"domain"` // e.g. "wastingnotime.org"
	Referrer    string            `json:"referrer"`
	ScreenWidth int               `json:"screen_width"`
	Props       map[string]string `json:"props"` // optional custom props
}

type QueueEvent struct {
	Domain      string            `json:"domain"`
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Referrer    string            `json:"referrer"`
	UserAgent   string            `json:"user_agent"`
	ScreenWidth int               `json:"screen_width"`
	IP          string            `json:"ip"`
	Props       map[string]string `json:"props"`
	Timestamp   string            `json:"timestamp"`
}

var (
	sqsClient *sqs.Client
	queueURL  string
)

func init() {
	queueURL = os.Getenv("EVENT_QUEUE_URL")
	if queueURL == "" {
		log.Fatal("EVENT_QUEUE_URL is required")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	sqsClient = sqs.NewFromConfig(cfg)
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var in PlausibleInput
	if err := json.Unmarshal([]byte(req.Body), &in); err != nil {
		log.Printf("invalid body: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error":"invalid body"}`,
		}, nil
	}

	ev := QueueEvent{
		Domain:      in.Domain,
		Name:        in.Name,
		URL:         in.URL,
		Referrer:    in.Referrer,
		UserAgent:   req.Headers["user-agent"],
		ScreenWidth: in.ScreenWidth,
		IP:          req.RequestContext.HTTP.SourceIP,
		Props:       in.Props,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	bodyBytes, err := json.Marshal(ev)
	if err != nil {
		log.Printf("marshal error: %v", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, nil
	}

	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(bodyBytes)),
	})
	if err != nil {
		log.Printf("failed to send to SQS: %v", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 202,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "https://wastingnotime.org", // CORS
			"Access-Control-Allow-Headers": "Content-Type",
		},
		Body: `{"status":"queued"}`,
	}, nil
}

func main() {
	lambda.Start(handler)
}
