package model

import (
	"github.com/google/uuid"
)

type Payment struct {
	ID          uuid.UUID `json:"id"`
	FromAccount int       `json:"from_account"`
	ToAccount   int       `json:"to_account"`
	Amount      int64     `json:"amount"`
	Status      string    `json:"status"`
}

type PaymentReq struct {
	FromAccount int    `json:"from_account"`
	ToAccount   int    `json:"to_account"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
}

type PaymentRes struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
