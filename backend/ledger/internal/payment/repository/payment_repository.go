package repository

import (
	"database/sql"
	"sort"

	"github.com/google/uuid"
	"github.com/rafaeldepontes/ledger/internal/payment"
	"github.com/rafaeldepontes/ledger/internal/payment/model"
	"github.com/rafaeldepontes/ledger/pkg/db/postgres"
	"github.com/rafaeldepontes/ledger/pkg/observability"
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
func (r repo) ProcessPayment(p model.Payment) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = lockAccounts(tx, p); err != nil {
		return err
	}

	if err = processPayment(tx, p); err != nil {
		return err
	}

	if err = processBalance(tx, p); err != nil {
		return err
	}

	return tx.Commit()
}

// GetPaymentByID implements [payment.Repository].
func (r repo) GetPaymentByID(id uuid.UUID) (model.Payment, error) {
	const query = `
	SELECT id, idempotency_key, from_account_id, to_account_id, amount, currency, status, created_at
	FROM payments
	WHERE id = $1
	`
	var p model.Payment
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.IdempotencyKey, &p.FromAccount, &p.ToAccount, &p.Amount, &p.Currency, &p.Status, &p.CreatedAt,
	)
	return p, err
}

// RefundPayment implements [payment.Repository].
func (r repo) RefundPayment(p model.Payment) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = lockAccounts(tx, p); err != nil {
		return err
	}

	if err = processRefund(tx, p); err != nil {
		return err
	}

	if err = processBalance(tx, p); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	observability.LedgerTransactionsTotal.Inc()
	return nil
}

// processPayment updates the Ledger table with the payment transaction.
func processPayment(tx *sql.Tx, p model.Payment) error {
	const paymentQuery = `
	INSERT INTO payments
	(id, idempotency_key, from_account_id, to_account_id, amount, currency, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.Exec(paymentQuery, p.ID, p.IdempotencyKey, p.FromAccount, p.ToAccount, p.Amount, p.Currency, p.Status)
	if err != nil {
		return err
	}

	const ledgerQuery = `INSERT INTO ledger_entries (payment_id, account_id, amount, currency) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(ledgerQuery, p.ID, p.FromAccount, -p.Amount, p.Currency)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ledgerQuery, p.ID, p.ToAccount, p.Amount, p.Currency)
	return err
}

// processBalance updates the Balance table for every new payment or refund made.
func processBalance(tx *sql.Tx, p model.Payment) error {
	const snapshotBalance = `
	INSERT INTO account_balance (account_id, balance, version)
	SELECT
		$1,
		COALESCE(SUM(amount), 0),
		(SELECT COALESCE(MAX(version), 0) + 1 FROM account_balance WHERE account_id = $1)
	FROM ledger_entries
	WHERE account_id = $1
	`
	_, err := tx.Exec(snapshotBalance, p.FromAccount)
	if err != nil {
		return err
	}

	_, err = tx.Exec(snapshotBalance, p.ToAccount)
	if err != nil {
		return err
	}

	return nil
}

// processRefund updates the Ledger and Payment tables with the new info.
func processRefund(tx *sql.Tx, p model.Payment) error {
	const updatePayment = `
	UPDATE payments
	SET status = 'refunded'
	WHERE id = $1
	`
	_, err := tx.Exec(updatePayment, p.ID)
	if err != nil {
		return err
	}

	const ledgerQuery = `INSERT INTO ledger_entries (payment_id, account_id, amount, currency) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(ledgerQuery, p.ID, p.FromAccount, p.Amount, p.Currency)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ledgerQuery, p.ID, p.ToAccount, -p.Amount, p.Currency)
	return err
}

// lockAccounts locks the account row that's going to be modified, it uses a sort algorithm
// internally.
func lockAccounts(tx *sql.Tx, p model.Payment) error {
	ids := []int{p.FromAccount, p.ToAccount}
	sort.Ints(ids)

	const lockAccounts = `SELECT id FROM accounts WHERE id = $1 OR id = $2 FOR UPDATE`
	_, err := tx.Exec(lockAccounts, ids[0], ids[1])
	return err
}
