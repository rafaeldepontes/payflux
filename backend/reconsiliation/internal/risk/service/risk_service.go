package service

import (
	"errors"
	"log"

	"github.com/google/uuid"
	reconciliation_model "github.com/rafaeldepontes/reconsiliation/internal/reconciliation/model"
	"github.com/rafaeldepontes/reconsiliation/internal/risk"
	"github.com/rafaeldepontes/reconsiliation/internal/risk/model"
	"github.com/rafaeldepontes/reconsiliation/pkg/observability"
)

type svc struct {
	repo risk.Repository
}

func NewService(repo risk.Repository) risk.Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) ProcessEvent(event reconciliation_model.PaymentEvent) error {
	log.Printf("[INFO] Evaluating risk for transaction: %s", event.PaymentID)

	paymentID, err := uuid.Parse(event.PaymentID)
	if err != nil {
		return err
	}

	score := 0
	flags := []string{}

	if event.Amount > 10000 {
		score += 50
		flags = append(flags, "LargeTransactionRule")
		observability.RiskFlagsTotal.Inc()
	}

	return s.repo.CreateRiskEvaluation(model.RiskEvaluation{
		TransactionID: paymentID,
		RiskScore:     score,
		Flags:         flags,
	})
}

func (s *svc) GetResult(txID string) (model.RiskEvaluation, error) {
	uid, err := uuid.Parse(txID)
	if err != nil {
		return model.RiskEvaluation{}, errors.New("invalid id")
	}
	return s.repo.GetRiskEvaluation(uid)
}
