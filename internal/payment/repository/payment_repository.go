package repository

import (
	"database/sql"

	"github.com/rafaeldepontes/goplo/internal/payment"
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
func (r repo) ProcessPayment(any) (string, error) {
	panic("unimplemented")
}
