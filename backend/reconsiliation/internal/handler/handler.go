package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/rafaeldepontes/reconsiliation/docs"
	"github.com/rafaeldepontes/reconsiliation/internal/rate/limit"
	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation"
	"github.com/rafaeldepontes/reconsiliation/internal/risk"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	ReconciliationC reconciliation.Controller
	RiskC           risk.Controller
}

func NewHandler(rc reconciliation.Controller, riskC risk.Controller, rateLimit limit.Middleware) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle(
		"GET /reconciliation/{transaction_id}", rateLimit.RateLimit(
			http.HandlerFunc(rc.GetReconciliationResult),
		),
	)

	mux.Handle(
		"GET /exceptions", rateLimit.RateLimit(
			http.HandlerFunc(rc.ListExceptions),
		),
	)

	mux.Handle(
		"POST /settlements", rateLimit.RateLimit(
			http.HandlerFunc(rc.CreateSettlementRecord),
		),
	)

	mux.Handle(
		"GET /risk/{transaction_id}", rateLimit.RateLimit(
			http.HandlerFunc(riskC.GetRiskEvaluation),
		),
	)

	// Observability
	mux.Handle("/metrics", promhttp.Handler())

	// Swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}
