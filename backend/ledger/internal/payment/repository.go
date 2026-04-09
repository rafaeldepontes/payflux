package payment

import (
	"github.com/google/uuid"
	"github.com/rafaeldepontes/ledger/internal/payment/model"
)

type Repository interface {
	ProcessPayment(p model.Payment) error
	GetPaymentByID(id uuid.UUID) (model.Payment, error)
	RefundPayment(p model.Payment) error
}
