package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/rafaeldepontes/reconsiliation/docs"
	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation"
	rs "github.com/rafaeldepontes/reconsiliation/internal/reconciliation/server"
	"github.com/rafaeldepontes/reconsiliation/internal/risk"
	risks "github.com/rafaeldepontes/reconsiliation/internal/risk/server"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	ReconciliationC reconciliation.Controller
	RiskC           risk.Controller
}

func newHandler() Handler {
	return Handler{
		ReconciliationC: rs.NewController(),
		RiskC:           risks.NewController(),
	}
}

func NewHandler() *http.ServeMux {
	handler := newHandler()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /reconciliation/{transaction_id}", handler.ReconciliationC.GetReconciliationResult)
	mux.HandleFunc("GET /exceptions", handler.ReconciliationC.ListExceptions)
	mux.HandleFunc("POST /settlements", handler.ReconciliationC.CreateSettlementRecord)

	mux.HandleFunc("GET /risk/{transaction_id}", handler.RiskC.GetRiskEvaluation)

	// Observability
	mux.Handle("/metrics", promhttp.Handler())

	// Swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}
