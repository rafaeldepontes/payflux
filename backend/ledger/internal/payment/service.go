package payment

import "github.com/rafaeldepontes/ledger/internal/payment/model"

type Service interface {
	ProcessPayment(key string, payment model.PaymentReq) (model.PaymentRes, error)
	CheckKey(key string) (model.PaymentRes, error)
	GetPayment(id string) (model.PaymentRes, error)
	RefundPayment(id string, req model.RefundReq) (model.PaymentRes, error)
}

type MessageBroker interface {
	Publish(body []byte) error
}
