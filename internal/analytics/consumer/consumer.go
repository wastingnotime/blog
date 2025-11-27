package consumer

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
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func ConfigFromEnv() (Config, error) {
	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		return Config{}, fmt.Errorf("SQS_QUEUE_URL not set")
	}

	plausibleURL := os.Getenv("PLAUSIBLE_URL")
	if plausibleURL == "" {
		plausibleURL = "http://localhost:8000/api/event"
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	if region == "" {
		region = "us-east-1" // fallback
	}

	return Config{
		QueueURL:     queueURL,
		PlausibleURL: plausibleURL,
		AWSRegion:    region,
	}, nil
}

func Run(ctx context.Context, cfg Config) error {
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.AWSRegion))
	if err != nil {
		return fmt.Errorf("unable to load AWS config: %w", err)
	}

	sqsClient := sqs.NewFromConfig(awsCfg)
	httpClient := &http.Client{Timeout: 10 * time.Second}

	for {
		select {
		case <-ctx.Done():
			log.Println("context cancelled, stopping consumer")
			return nil
		default:
		}

		out, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(cfg.QueueURL),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     20,
		})
		if err != nil {
			log.Printf("error receiving messages: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if len(out.Messages) == 0 {
			continue
		}

		for _, m := range out.Messages {
			if m.Body == nil {
				continue
			}

			if err := processMessage(ctx, httpClient, sqsClient, cfg, &m); err != nil {
				log.Printf("failed processing message: %v", err)
				// do NOT delete message; let SQS redrive / retry
				continue
			}

			_, err := sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(cfg.QueueURL),
				ReceiptHandle: m.ReceiptHandle,
			})
			if err != nil {
				log.Printf("failed to delete message: %v", err)
			}
		}
	}
}

func processMessage(
	ctx context.Context,
	httpClient *http.Client,
	sqsClient *sqs.Client,
	cfg Config,
	msg *types.Message,
) error {
	var qev QueueEvent
	if err := json.Unmarshal([]byte(*msg.Body), &qev); err != nil {
		return fmt.Errorf("unmarshal queue event: %w", err)
	}

	pev := PlausibleEvent{
		Name:        qev.Name,
		URL:         qev.URL,
		Domain:      qev.Domain,
		Referrer:    qev.Referrer,
		ScreenWidth: qev.ScreenWidth,
		Props:       qev.Props,
		Timestamp:   qev.Timestamp,
	}

	body, err := json.Marshal(pev)
	if err != nil {
		return fmt.Errorf("marshal plausible event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.PlausibleURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if qev.UserAgent != "" {
		req.Header.Set("User-Agent", qev.UserAgent)
	}
	if qev.IP != "" {
		// Plausible uses X-Forwarded-For to infer client IP
		req.Header.Set("X-Forwarded-For", qev.IP)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("post to plausible: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("plausible returned status %d", resp.StatusCode)
	}

	return nil
}
