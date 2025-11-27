package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type PlausibleEvent struct {
	SiteID      string `json:"site_id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Referrer    string `json:"referrer"`
	UserAgent   string `json:"user_agent"`
	ScreenWidth int    `json:"screen_width"`
	Timestamp   string `json:"timestamp"`
	IP          string `json:"ip"`
}

func main() {
	ctx := context.Background()

	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		log.Fatal("SQS_QUEUE_URL not set")
	}

	plausibleURL := os.Getenv("PLAUSIBLE_URL")
	if plausibleURL == "" {
		plausibleURL = "http://localhost:8000/api/event"
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	log.Println("Starting SQS consumer...")
	for {
		output, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     20, // long polling
		})
		if err != nil {
			log.Printf("error receiving messages: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if len(output.Messages) == 0 {
			continue
		}

		for _, m := range output.Messages {
			if m.Body == nil {
				continue
			}

			var evt PlausibleEvent
			if err := json.Unmarshal([]byte(*m.Body), &evt); err != nil {
				log.Printf("invalid message body: %v", err)
				// maybe move to DLQ in real scenario
				continue
			}

			if err := sendToPlausible(plausibleURL, evt); err != nil {
				log.Printf("failed to send to Plausible: %v", err)
				// don't delete message, so it can be retried later
				continue
			}

			// delete only after successful processing
			_, err := sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: m.ReceiptHandle,
			})
			if err != nil {
				log.Printf("failed to delete message: %v", err)
			}
		}
	}
}

func sendToPlausible(url string, evt PlausibleEvent) error {
	body, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", evt.UserAgent)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("plausible returned status %d", resp.StatusCode)
	}

	return nil
}
