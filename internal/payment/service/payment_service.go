package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rafaeldepontes/goplo/internal/cache"
	"github.com/rafaeldepontes/goplo/internal/payment"

	cs "github.com/rafaeldepontes/goplo/internal/cache/service"
	pr "github.com/rafaeldepontes/goplo/internal/payment/repository"
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
func (s service) ProcessPayment(key string) (string, error) {
	paymentID := uuid.New().String()

	// TODO: we will also save to the repository here.
	// _, err := s.repository.ProcessPayment(nil)
	// if err != nil {
	// 	return "", err
	// }

	s.cache.Add(key, paymentID)
	return paymentID, nil
}

// CheckKey checks if the idempotency key is already in the cache.
func (s service) CheckKey(key string) (string, error) {
	val, has := s.cache.Get(key)
	if !has {
		return "", errors.New("not on cache")
	}
	return val, nil
}
