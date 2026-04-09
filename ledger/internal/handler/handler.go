package handler

import (
	"net/http"

	"github.com/rafaeldepontes/goplo/internal/account"
	as "github.com/rafaeldepontes/goplo/internal/account/server"
	"github.com/rafaeldepontes/goplo/internal/payment"
	ps "github.com/rafaeldepontes/goplo/internal/payment/server"
)

type Handler struct {
	// Controllers...
	PaymentC payment.Controller
	AccountC account.Controller
}

func newHandler() Handler {
	return Handler{
		PaymentC: ps.NewController(),
		AccountC: as.NewController(),
	}
}

func NewHandler() *http.ServeMux {
	handler := newHandler()
	mux := http.NewServeMux()

	// Payment Routes
	mux.HandleFunc("POST /payments", handler.PaymentC.ProcessPayment)
	mux.HandleFunc("GET /payments/{id}", handler.PaymentC.GetPayment)
	mux.HandleFunc("POST /payments/{id}/refund", handler.PaymentC.RefundPayment)

	// Account Routes
	mux.HandleFunc("GET /accounts/{id}/balance", handler.AccountC.GetAccountBalance)

	return mux
}
