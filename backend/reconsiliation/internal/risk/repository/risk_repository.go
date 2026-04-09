package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rafaeldepontes/reconsiliation/internal/risk"
	"github.com/rafaeldepontes/reconsiliation/internal/risk/model"
	"github.com/rafaeldepontes/reconsiliation/pkg/db/postgres"
)

type repo struct {
	db *sql.DB
}

func NewRepository() risk.Repository {
	return &repo{
		db: postgres.GetDb(),
	}
}

func (r *repo) CreateRiskEvaluation(evaluation model.RiskEvaluation) error {
	const query = `
	INSERT INTO risk_evaluations (transaction_id, risk_score, flags)
	VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, evaluation.TransactionID, evaluation.RiskScore, pq.Array(evaluation.Flags))
	return err
}

func (r *repo) GetRiskEvaluation(txID uuid.UUID) (model.RiskEvaluation, error) {
	const query = `
	SELECT id, transaction_id, risk_score, flags, created_at
	FROM risk_evaluations
	WHERE transaction_id = $1
	ORDER BY created_at DESC LIMIT 1
	`
	var res model.RiskEvaluation
	err := r.db.QueryRow(query, txID).Scan(&res.ID, &res.TransactionID, &res.RiskScore, pq.Array(&res.Flags), &res.CreatedAt)
	return res, err
}
