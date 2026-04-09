package model

import (
	"time"

	"github.com/google/uuid"
)

type Exception struct {
	ID               uuid.UUID `json:"id"`
	TransactionID    uuid.UUID `json:"transaction_id"`
	Type             string    `json:"type"`
	LedgerAmount     int64     `json:"ledger_amount"`
	SettlementAmount int64     `json:"settlement_amount"`
	CreatedAt        time.Time `json:"created_at"`
}
