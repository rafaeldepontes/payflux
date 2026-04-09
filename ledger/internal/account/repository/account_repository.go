package repository

import (
	"database/sql"

	"github.com/rafaeldepontes/goplo/internal/account"
	"github.com/rafaeldepontes/goplo/pkg/db/postgres"
)

type repo struct {
	db *sql.DB
}

func NewRepository() account.Repository {
	return repo{
		db: postgres.GetDb(),
	}
}

func (r repo) GetAccountBalance(accountID int) (int64, error) {
	const query = `
	SELECT COALESCE(SUM(amount), 0)
	FROM ledger_entries
	WHERE account_id = $1
	`
	var balance int64
	err := r.db.QueryRow(query, accountID).Scan(&balance)
	return balance, err
}
