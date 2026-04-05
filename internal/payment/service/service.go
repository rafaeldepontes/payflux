package service

import "github.com/rafaeldepontes/goplo/internal/payment"

type service struct {
	// ...
}

func NewService() payment.Service {
	return service{}
}

func (s service) ProcessPayment() (string, error) {
	return "", nil
}
