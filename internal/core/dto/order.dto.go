package dto

type WebhookEventDTO struct {
	ID          string
	LiveMode    bool
	Type        string
	DateCreated string
	UserID      any
	APIVersion  string
	Action      string
	OrderID     string
}
