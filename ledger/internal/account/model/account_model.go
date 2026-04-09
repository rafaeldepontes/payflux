package model

import "time"

type Account struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type BalanceRes struct {
	AccountID int   `json:"account_id"`
	Balance   int64 `json:"balance"`
}
