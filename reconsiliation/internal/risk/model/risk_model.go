package model

import (
	"time"

	"github.com/google/uuid"
)

type RiskEvaluation struct {
	ID            uuid.UUID `json:"id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	RiskScore     int       `json:"risk_score"`
	Flags         []string  `json:"flags"`
	CreatedAt     time.Time `json:"created_at"`
}
