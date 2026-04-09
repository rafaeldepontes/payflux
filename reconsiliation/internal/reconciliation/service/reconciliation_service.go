package service

import (
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation"
	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation/model"
	rr "github.com/rafaeldepontes/reconsiliation/internal/reconciliation/repository"
)

type svc struct {
	repo reconciliation.Repository
}

func NewService() reconciliation.Service {
	return &svc{
		repo: rr.NewRepository(),
	}
}

func (s *svc) ProcessEvent(event model.PaymentEvent) error {
	log.Printf("[INFO] Reconciling transaction: %s", event.PaymentID)

	paymentID, err := uuid.Parse(event.PaymentID)
	if err != nil {
		return err
	}

	settlement, err := s.repo.GetSettlementRecord(paymentID)
	if err == sql.ErrNoRows {
		log.Printf("[WARN] Missing settlement record for transaction: %s", event.PaymentID)
		return s.repo.CreateException(model.Exception{
			TransactionID:    paymentID,
			Type:             "MissingSettlementRecord",
			LedgerAmount:     event.Amount,
			SettlementAmount: 0,
		})
	} else if err != nil {
		return err
	}

	status := "matched"
	if settlement.Amount != event.Amount {
		status = "mismatched"
		_ = s.repo.CreateException(model.Exception{
			TransactionID:    paymentID,
			Type:             "AmountMismatch",
			LedgerAmount:     event.Amount,
			SettlementAmount: settlement.Amount,
		})
	}

	return s.repo.CreateReconciliationResult(model.ReconciliationResult{
		TransactionID:    paymentID,
		Status:           status,
		LedgerAmount:     event.Amount,
		SettlementAmount: settlement.Amount,
	})
}

func (s *svc) GetResult(txID string) (model.ReconciliationResult, error) {
	uid, err := uuid.Parse(txID)
	if err != nil {
		return model.ReconciliationResult{}, errors.New("invalid id")
	}
	return s.repo.GetReconciliationResult(uid)
}

func (s *svc) ListExceptions() ([]model.Exception, error) {
	return s.repo.ListExceptions()
}

func (s *svc) CreateSettlementRecord(txID string, amount int64, status string) error {
	uid, err := uuid.Parse(txID)
	if err != nil {
		return errors.New("invalid id")
	}
	return s.repo.CreateSettlementRecord(model.SettlementRecord{
		TransactionID: uid,
		Amount:        amount,
		Status:        status,
	})
}
