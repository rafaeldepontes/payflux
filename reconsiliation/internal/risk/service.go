package risk

import (
	reconciliation_model "github.com/rafaeldepontes/reconsiliation/internal/reconciliation/model"
	"github.com/rafaeldepontes/reconsiliation/internal/risk/model"
)

type Service interface {
	ProcessEvent(event reconciliation_model.PaymentEvent) error
	GetResult(txID string) (model.RiskEvaluation, error)
}
