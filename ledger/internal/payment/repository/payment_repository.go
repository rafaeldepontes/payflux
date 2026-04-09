package repository

import (
	"database/sql"

	"github.com/rafaeldepontes/goplo/internal/payment"
	"github.com/rafaeldepontes/goplo/internal/payment/model"
	"github.com/rafaeldepontes/goplo/pkg/db/postgres"
)

type repo struct {
	db *sql.DB
}

func NewRepository() payment.Repository {
	return repo{
		db: postgres.GetDb(),
	}
}

// ProcessPayment implements [payment.Repository].
func (r repo) ProcessPayment(p model.Payment, key, currency string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const paymentQuery = `
	INSERT INTO payments 
	(id, idempotency_key, from_account_id, to_account_id, amount, currency, status) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(paymentQuery, p.ID, key, p.FromAccount, p.ToAccount, p.Amount, currency, p.Status)
	if err != nil {
		return err
	}

	const ledgerQuery = `INSERT INTO ledger_entries (payment_id, account_id, amount, currency) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(ledgerQuery, p.ID, p.FromAccount, -p.Amount, currency)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ledgerQuery, p.ID, p.ToAccount, p.Amount, currency)
	if err != nil {
		return err
	}

	return tx.Commit()
}
