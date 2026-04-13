package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/rafaeldepontes/ledger/docs"
	"github.com/rafaeldepontes/ledger/internal/account"
	"github.com/rafaeldepontes/ledger/internal/payment"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	PaymentC payment.Controller
	AccountC account.Controller
}

func NewHandler(paymentC payment.Controller, accountC account.Controller) *http.ServeMux {
	mux := http.NewServeMux()

	// Payment Routes
	mux.HandleFunc("POST /payments", paymentC.ProcessPayment)
	mux.HandleFunc("GET /payments/{id}", paymentC.GetPayment)
	mux.HandleFunc("POST /payments/{id}/refund", paymentC.RefundPayment)

	// Account Routes
	mux.HandleFunc("GET /accounts/{id}/balance", accountC.GetAccountBalance)

	// Observability
	mux.Handle("/metrics", promhttp.Handler())

	// Swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}
