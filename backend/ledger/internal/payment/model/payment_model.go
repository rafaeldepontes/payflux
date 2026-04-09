package model

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID             uuid.UUID `json:"id" db:"id"`
	IdempotencyKey string    `json:"idempotency_key" db:"idempotency_key"`
	FromAccount    int       `json:"from_account" db:"from_account_id"`
	ToAccount      int       `json:"to_account" db:"to_account_id"`
	Amount         int64     `json:"amount" db:"amount"`
	Status         string    `json:"status" db:"status"`
	Currency       string    `json:"currency" db:"currency"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type PaymentReq struct {
	FromAccount int    `json:"from_account"`
	ToAccount   int    `json:"to_account"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
}

type PaymentRes struct {
	ID       string `json:"payment_id"`
	Status   string `json:"status"`
	Amount   int64  `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type RefundReq struct {
	Amount int64 `json:"amount"`
}

type PaymentEvent struct {
	EventType string    `json:"event_type"`
	PaymentID string    `json:"payment_id"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}
