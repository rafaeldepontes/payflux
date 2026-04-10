package service

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rafaeldepontes/ledger/internal/cache"
	"github.com/rafaeldepontes/ledger/internal/payment"

	cs "github.com/rafaeldepontes/ledger/internal/cache/service"
	"github.com/rafaeldepontes/ledger/internal/payment/model"
	pr "github.com/rafaeldepontes/ledger/internal/payment/repository"
	"github.com/rafaeldepontes/ledger/pkg/message-broker/rabbitmq"
	"github.com/rafaeldepontes/ledger/pkg/observability"
)

const (
	CompletedStatus = "completed"
	RefundedStatus  = "refunded"
)

type service struct {
	cache      cache.Cache[string, string]
	repository payment.Repository
}

// NewService returns a new instance of the payment service.
func NewService() payment.Service {
	return service{
		cache:      cs.NewService(),
		repository: pr.NewRepository(),
	}
}

// ProcessPayment generates a unique payment ID and stores it in the cache with the idempotency key.
func (s service) ProcessPayment(key string, payment model.PaymentReq) (model.PaymentRes, error) {
	if err := validatePaymentRequest(payment); err != nil {
		return model.PaymentRes{}, err
	}

	p := model.Payment{
		ID:             uuid.New(),
		IdempotencyKey: key,
		FromAccount:    payment.FromAccount,
		ToAccount:      payment.ToAccount,
		Amount:         payment.Amount,
		Status:         CompletedStatus,
		Currency:       payment.Currency,
	}

	err := s.repository.ProcessPayment(p)
	if err != nil {
		observability.PaymentFailuresTotal.Inc()
		log.Println("[ERROR] could not finish the payment:", err)
		return model.PaymentRes{}, errors.New("something went wrong")
	}

	event := model.PaymentEvent{
		EventType: "PaymentCompleted",
		PaymentID: p.ID.String(),
		Amount:    p.Amount,
		Currency:  p.Currency,
		Timestamp: time.Now(),
	}
	eventBody, _ := json.Marshal(event)
	if err := rabbitmq.Publish(eventBody); err != nil {
		log.Println("[WARN] could not publish event:", err)
	}

	res := model.PaymentRes{
		ID:       p.ID.String(),
		Status:   p.Status,
		Amount:   p.Amount,
		Currency: p.Currency,
	}

	body, _ := json.Marshal(res)
	s.cache.Add(key, string(body))
	return res, nil
}

// CheckKey checks if the idempotency key is already in the cache.
func (s service) CheckKey(key string) (model.PaymentRes, error) {
	if key == "" {
		return model.PaymentRes{}, errors.New("not on cache")
	}

	val, has := s.cache.Get(key)
	if !has {
		return model.PaymentRes{}, errors.New("not on cache")
	}

	var res model.PaymentRes
	err := json.Unmarshal([]byte(val), &res)
	return res, err
}

// GetPayment retrieves a payment by its ID.
func (s service) GetPayment(id string) (model.PaymentRes, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return model.PaymentRes{}, errors.New("invalid id")
	}

	p, err := s.repository.GetPaymentByID(uid)
	if err != nil {
		return model.PaymentRes{}, errors.New("payment not found")
	}

	return model.PaymentRes{
		ID:       p.ID.String(),
		Status:   p.Status,
		Amount:   p.Amount,
		Currency: p.Currency,
	}, nil
}

// RefundPayment processes a refund for a given payment.
func (s service) RefundPayment(id string, req model.RefundReq) (model.PaymentRes, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return model.PaymentRes{}, errors.New("invalid id")
	}

	p, err := s.repository.GetPaymentByID(uid)
	if err != nil {
		return model.PaymentRes{}, errors.New("payment not found")
	}

	if p.Status == RefundedStatus {
		return model.PaymentRes{}, errors.New("payment already refunded")
	}

	// For simplicity, we only allow full refund if amount matches or if no amount provided in req
	if req.Amount > 0 && req.Amount != p.Amount {
		return model.PaymentRes{}, errors.New("partial refund not supported yet")
	}

	err = s.repository.RefundPayment(p)
	if err != nil {
		return model.PaymentRes{}, errors.New("could not process refund")
	}

	// Emit event to rabbitmq
	event := model.PaymentEvent{
		EventType: "PaymentRefunded",
		PaymentID: p.ID.String(),
		Amount:    p.Amount,
		Currency:  p.Currency,
		Timestamp: time.Now(),
	}
	eventBody, _ := json.Marshal(event)
	if err := rabbitmq.Publish(eventBody); err != nil {
		log.Println("[WARN] could not publish refund event:", err)
	}

	return model.PaymentRes{
		ID:     p.ID.String(),
		Status: RefundedStatus,
	}, nil
}

func validatePaymentRequest(p model.PaymentReq) error {
	if p.FromAccount <= 0 {
		return errors.New("source account is required")
	}

	if p.ToAccount <= 0 {
		return errors.New("destination account is required")
	}

	if p.Amount <= 0 {
		return errors.New("amount needs to be greater than zero")
	}

	if strings.TrimSpace(p.Currency) == "" {
		return errors.New("currency is required")
	}

	return nil
}
