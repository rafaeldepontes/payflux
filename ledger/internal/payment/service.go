package payment

import "github.com/rafaeldepontes/goplo/internal/payment/model"

type Service interface {
	ProcessPayment(key string, payment model.PaymentReq) (model.PaymentRes, error)
	CheckKey(key string) (model.PaymentRes, error)
}
