package service

import (
	"errors"

	"github.com/rafaeldepontes/goplo/internal/cache"
	"github.com/rafaeldepontes/goplo/internal/payment"

	cs "github.com/rafaeldepontes/goplo/internal/cache/service"
	pr "github.com/rafaeldepontes/goplo/internal/payment/repository"
)

type service struct {
	cache      cache.Cache[string, string]
	repository payment.Repository
}

func NewService() payment.Service {
	return service{
		cache:      cs.NewService(),
		repository: pr.NewRepository(),
	}
}

// TODO: finish the logic behind the double ledger payment process
func (s service) ProcessPayment() (string, error) {
	return "", nil
}

func (s service) CheckKey(key string) (string, error) {
	val, has := s.cache.Get(key)
	if !has {
		return "", errors.New("not on cache")
	}
	return val, nil
}
