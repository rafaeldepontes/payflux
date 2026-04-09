package reconciliation

import (
	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation/model"
)

type Service interface {
	ProcessEvent(event model.PaymentEvent) error
	GetResult(txID string) (model.ReconciliationResult, error)
	ListExceptions() ([]model.Exception, error)
	CreateSettlementRecord(txID string, amount int64, status string) error
}
