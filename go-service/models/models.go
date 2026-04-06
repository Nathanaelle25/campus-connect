package models

// WebhookPayload represents the structure sent by NestJS
type WebhookPayload struct {
	EventID   int    `json:"eventId"`
	Title     string `json:"title"`
	Action    string `json:"action"`
	Timestamp string `json:"timestamp"`
}
