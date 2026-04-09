package payment

import "github.com/rafaeldepontes/goplo/internal/payment/model"

type Repository interface {
	ProcessPayment(p model.Payment, key, currency string) error
}
