package eventsink

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// InboundEvent matches the JSON payload received from clients.
type InboundEvent struct {
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Domain      string            `json:"domain"`
	Referrer    string            `json:"referrer"`
	ScreenWidth int               `json:"screen_width"`
	Props       map[string]string `json:"props"`
}

// QueueEvent is the message body sent to SQS for downstream processing.
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

// Meta holds request metadata that is not part of the client payload.
type Meta struct {
	UserAgent string
	IP        string
}

// Processor validates inbound events and enqueues them to SQS.
type Processor struct {
	sqsClient *sqs.Client
	queueURL  string
	now       func() time.Time
}

// ErrInvalidEvent indicates a missing required field.
var ErrInvalidEvent = errors.New("invalid event")

// NewProcessor constructs a Processor with defaults.
func NewProcessor(sqsClient *sqs.Client, queueURL string) *Processor {
	return &Processor{
		sqsClient: sqsClient,
		queueURL:  queueURL,
		now:       time.Now,
	}
}

// Handle validates and sends an inbound event to the queue.
func (p *Processor) Handle(ctx context.Context, in InboundEvent, meta Meta) error {
	if in.Name == "" || in.URL == "" || in.Domain == "" {
		return ErrInvalidEvent
	}

	if in.Props == nil {
		in.Props = make(map[string]string)
	}

	qev := QueueEvent{
		Domain:      in.Domain,
		Name:        in.Name,
		URL:         in.URL,
		Referrer:    in.Referrer,
		UserAgent:   meta.UserAgent,
		ScreenWidth: in.ScreenWidth,
		IP:          meta.IP,
		Timestamp:   p.now().UTC().Format(time.RFC3339),
		Props:       in.Props,
	}

	body, err := json.Marshal(qev)
	if err != nil {
		return err
	}

	_, err = p.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(string(body)),
	})
	return err
}
