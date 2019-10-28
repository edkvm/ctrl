package trigger

type WebhookID string

type Webhook struct {
	ID        WebhookID
	Action    string
	Enabled   bool
}
