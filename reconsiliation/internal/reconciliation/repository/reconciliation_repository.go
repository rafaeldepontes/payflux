package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation"
	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation/model"
	"github.com/rafaeldepontes/reconsiliation/pkg/db/postgres"
)

type repo struct {
	db *sql.DB
}

func NewRepository() reconciliation.Repository {
	return &repo{
		db: postgres.GetDb(),
	}
}

func (r *repo) GetSettlementRecord(txID uuid.UUID) (model.SettlementRecord, error) {
	const query = `SELECT id, transaction_id, amount, status FROM settlement_records WHERE transaction_id = $1`
	var settlement model.SettlementRecord
	err := r.db.QueryRow(query, txID).Scan(&settlement.ID, &settlement.TransactionID, &settlement.Amount, &settlement.Status)
	return settlement, err
}

func (r *repo) CreateSettlementRecord(record model.SettlementRecord) error {
	const query = `
	INSERT INTO settlement_records (transaction_id, amount, status)
	VALUES ($1, $2, $3)
	ON CONFLICT (transaction_id) DO UPDATE SET amount = $2, status = $3
	`
	_, err := r.db.Exec(query, record.TransactionID, record.Amount, record.Status)
	return err
}

func (r *repo) CreateReconciliationResult(res model.ReconciliationResult) error {
	const query = `
	INSERT INTO reconciliation_results (transaction_id, status, ledger_amount, settlement_amount)
	VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, res.TransactionID, res.Status, res.LedgerAmount, res.SettlementAmount)
	return err
}

func (r *repo) GetReconciliationResult(txID uuid.UUID) (model.ReconciliationResult, error) {
	const query = `
	SELECT id, transaction_id, status, ledger_amount, settlement_amount, created_at
	FROM reconciliation_results
	WHERE transaction_id = $1
	ORDER BY created_at DESC LIMIT 1
	`
	var res model.ReconciliationResult
	err := r.db.QueryRow(query, txID).Scan(&res.ID, &res.TransactionID, &res.Status, &res.LedgerAmount, &res.SettlementAmount, &res.CreatedAt)
	return res, err
}

func (r *repo) CreateException(exc model.Exception) error {
	const query = `
	INSERT INTO exceptions (transaction_id, type, ledger_amount, settlement_amount)
	VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, exc.TransactionID, exc.Type, exc.LedgerAmount, exc.SettlementAmount)
	return err
}

func (r *repo) ListExceptions() ([]model.Exception, error) {
	const query = `
	SELECT id, transaction_id, type, ledger_amount, settlement_amount, created_at
	FROM exceptions
	ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exceptions []model.Exception
	for rows.Next() {
		var e model.Exception
		if err := rows.Scan(&e.ID, &e.TransactionID, &e.Type, &e.LedgerAmount, &e.SettlementAmount, &e.CreatedAt); err != nil {
			return nil, err
		}
		exceptions = append(exceptions, e)
	}
	return exceptions, nil
}
