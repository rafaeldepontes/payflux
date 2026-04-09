package handler

import (
	"net/http"

	"github.com/rafaeldepontes/goplo/internal/payment"
	ps "github.com/rafaeldepontes/goplo/internal/payment/server"
)

type Handler struct {
	// Controllers...
	PaymentC payment.Controller
}

func newHandler() Handler {
	return Handler{
		PaymentC: ps.NewController(),
	}
}

func NewHandler() *http.ServeMux {
	handler := newHandler()
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("POST /payments", handler.PaymentC.ProcessPayment)

	return mux
}
