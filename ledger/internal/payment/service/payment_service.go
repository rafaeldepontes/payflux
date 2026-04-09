package service

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/rafaeldepontes/goplo/internal/cache"
	"github.com/rafaeldepontes/goplo/internal/payment"

	cs "github.com/rafaeldepontes/goplo/internal/cache/service"
	"github.com/rafaeldepontes/goplo/internal/payment/model"
	pr "github.com/rafaeldepontes/goplo/internal/payment/repository"
)

const (
	CompletedStatus = "completed"
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
		log.Println("[ERROR] could not finish the payment:", err)
		return model.PaymentRes{}, errors.New("something went wrong")
	}

	// TODO: implement rabbitmq producer for analytics...
	// ...

	res := model.PaymentRes{
		ID:     p.ID.String(),
		Status: p.Status,
	}

	body, _ := json.Marshal(res)
	s.cache.Add(key, string(body))
	return res, nil
}

// CheckKey checks if the idempotency key is already in the cache.
func (s service) CheckKey(key string) (model.PaymentRes, error) {
	val, has := s.cache.Get(key)
	if !has {
		return model.PaymentRes{}, errors.New("not on cache")
	}

	var res model.PaymentRes
	err := json.Unmarshal([]byte(val), &res)
	return res, err
}
