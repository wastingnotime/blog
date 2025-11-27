package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type InboundEvent struct {
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Domain      string            `json:"domain"`
	Referrer    string            `json:"referrer"`
	ScreenWidth int               `json:"screen_width"`
	Props       map[string]string `json:"props"`
}

type QueueEvent struct {
	Domain      string            `json:"domain"`
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Referrer    string            `json:"referrer"`
	UserAgent   string            `json:"user_agent"`
	ScreenWidth int               `json:"screen_width"`
	IP          string            `json:"ip"`
	Timestamp   string            `json:"timestamp"` // RFC3339
	Props       map[string]string `json:"props"`
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
	log.Printf("Lambda initialized. Queue URL: %s", queueURL)
}

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Handle CORS preflight if it reaches Lambda (often API GW does this itself)
	if strings.ToUpper(req.RequestContext.HTTP.Method) == http.MethodOptions {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusNoContent,
			Headers:    corsHeaders(),
		}, nil
	}

	if req.Body == "" {
		return errorResponse(http.StatusBadRequest, "empty body"), nil
	}

	var in InboundEvent
	if err := json.Unmarshal([]byte(req.Body), &in); err != nil {
		log.Printf("invalid JSON body: %v", err)
		return errorResponse(http.StatusBadRequest, "invalid JSON"), nil
	}

	// Basic validation
	if in.Name == "" || in.URL == "" || in.Domain == "" {
		return errorResponse(http.StatusBadRequest, "name, url and domain are required"), nil
	}

	ua := req.Headers["user-agent"]
	ip := extractIP(req)

	if in.Props == nil {
		in.Props = make(map[string]string)
	}

	ev := QueueEvent{
		Domain:      in.Domain,
		Name:        in.Name,
		URL:         in.URL,
		Referrer:    in.Referrer,
		UserAgent:   ua,
		ScreenWidth: in.ScreenWidth,
		IP:          ip,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Props:       in.Props,
	}

	bodyBytes, err := json.Marshal(ev)
	if err != nil {
		log.Printf("failed to marshal queue event: %v", err)
		return errorResponse(http.StatusInternalServerError, "internal error"), nil
	}

	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(bodyBytes)),
	})
	if err != nil {
		log.Printf("failed to send message to SQS: %v", err)
		return errorResponse(http.StatusInternalServerError, "failed to queue event"), nil
	}

	respBody, _ := json.Marshal(map[string]string{
		"status": "queued",
	})

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusAccepted,
		Body:       string(respBody),
		Headers:    corsHeaders(),
	}, nil
}

func extractIP(req events.APIGatewayV2HTTPRequest) string {
	// Prefer X-Forwarded-For if present
	if xff, ok := req.Headers["x-forwarded-for"]; ok && xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	// Fallback to SourceIP from request context
	if req.RequestContext.HTTP.SourceIP != "" {
		return req.RequestContext.HTTP.SourceIP
	}

	return ""
}

func corsHeaders() map[string]string {
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		// fallback – you can tighten this in prod
		allowedOrigin = "*"
	}

	return map[string]string{
		"Access-Control-Allow-Origin":      allowedOrigin,
		"Access-Control-Allow-Methods":     "OPTIONS,POST",
		"Access-Control-Allow-Headers":     "Content-Type",
		"Access-Control-Allow-Credentials": "false",
	}
}

func errorResponse(status int, msg string) events.APIGatewayV2HTTPResponse {
	body, _ := json.Marshal(map[string]string{
		"error": msg,
	})

	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Body:       string(body),
		Headers:    corsHeaders(),
	}
}

func main() {
	lambda.Start(handler)
}
