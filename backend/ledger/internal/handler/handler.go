package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/rafaeldepontes/ledger/docs"
	"github.com/rafaeldepontes/ledger/internal/account"
	as "github.com/rafaeldepontes/ledger/internal/account/server"
	"github.com/rafaeldepontes/ledger/internal/payment"
	ps "github.com/rafaeldepontes/ledger/internal/payment/server"
	httpSwagger "github.com/swaggo/http-swagger"
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

	// Observability
	mux.Handle("/metrics", promhttp.Handler())

	// Swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}
