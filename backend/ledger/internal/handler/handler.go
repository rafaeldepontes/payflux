package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/rafaeldepontes/ledger/docs"
	"github.com/rafaeldepontes/ledger/internal/account"
	"github.com/rafaeldepontes/ledger/internal/payment"
	"github.com/rafaeldepontes/ledger/internal/rate/limit"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	PaymentC payment.Controller
	AccountC account.Controller
}

func NewHandler(paymentC payment.Controller, accountC account.Controller, rateLimit limit.Middleware) *http.ServeMux {
	mux := http.NewServeMux()

	// Payment Routes
	mux.Handle("POST /payments", rateLimit.RateLimit(http.HandlerFunc(paymentC.ProcessPayment)))
	mux.Handle("GET /payments/{id}", rateLimit.RateLimit(http.HandlerFunc(paymentC.GetPayment)))
	mux.Handle("POST /payments/{id}/refund", rateLimit.RateLimit(http.HandlerFunc(paymentC.RefundPayment)))

	// Account Routes
	mux.Handle("GET /accounts/{id}/balance", rateLimit.RateLimit(http.HandlerFunc(accountC.GetAccountBalance)))

	// Observability
	mux.Handle("/metrics", promhttp.Handler())

	// Swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}
