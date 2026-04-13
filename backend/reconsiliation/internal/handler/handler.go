package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/rafaeldepontes/reconsiliation/docs"
	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation"
	"github.com/rafaeldepontes/reconsiliation/internal/risk"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	ReconciliationC reconciliation.Controller
	RiskC           risk.Controller
}

func NewHandler(reconciliationC reconciliation.Controller, riskC risk.Controller) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /reconciliation/{transaction_id}", reconciliationC.GetReconciliationResult)
	mux.HandleFunc("GET /exceptions", reconciliationC.ListExceptions)
	mux.HandleFunc("POST /settlements", reconciliationC.CreateSettlementRecord)

	mux.HandleFunc("GET /risk/{transaction_id}", riskC.GetRiskEvaluation)

	// Observability
	mux.Handle("/metrics", promhttp.Handler())

	// Swagger
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	return mux
}
