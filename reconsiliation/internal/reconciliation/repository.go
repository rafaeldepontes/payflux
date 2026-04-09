package reconciliation

import (
	"github.com/google/uuid"
	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation/model"
)

type Repository interface {
	GetSettlementRecord(txID uuid.UUID) (model.SettlementRecord, error)
	CreateSettlementRecord(record model.SettlementRecord) error
	CreateReconciliationResult(res model.ReconciliationResult) error
	GetReconciliationResult(txID uuid.UUID) (model.ReconciliationResult, error)
	CreateException(exc model.Exception) error
	ListExceptions() ([]model.Exception, error)
}
