package risk

import (
	"github.com/google/uuid"
	"github.com/rafaeldepontes/reconsiliation/internal/risk/model"
)

type Repository interface {
	CreateRiskEvaluation(evaluation model.RiskEvaluation) error
	GetRiskEvaluation(txID uuid.UUID) (model.RiskEvaluation, error)
}
