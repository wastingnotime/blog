package consumer

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

// PlausibleEvent - payload expected by Plausible's /api/event (body)
type PlausibleEvent struct {
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Domain      string            `json:"domain"`
	Referrer    string            `json:"referrer,omitempty"`
	ScreenWidth int               `json:"screen_width,omitempty"`
	Props       map[string]string `json:"props,omitempty"`
	Timestamp   string            `json:"timestamp,omitempty"` // RFC3339
}

// Config - runtime config for the consumer
type Config struct {
	QueueURL     string
	PlausibleURL string
	AWSRegion    string
}
