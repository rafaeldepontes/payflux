package model

import (
	"time"

	"github.com/google/uuid"
)

type LedgerEntry struct {
	ID        uuid.UUID `json:"id"`
	PaymentID uuid.UUID `json:"payment_id"`
	AccountID int       `json:"account_id"`
	Amount    int64     `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}
