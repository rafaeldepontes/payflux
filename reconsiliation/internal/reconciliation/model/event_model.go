package model

import (
	"time"
)

type PaymentEvent struct {
	EventType string    `json:"event_type"`
	PaymentID string    `json:"payment_id"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}
