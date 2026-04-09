package handler

import (
	"net/http"

	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation"
	rs "github.com/rafaeldepontes/reconsiliation/internal/reconciliation/server"
	"github.com/rafaeldepontes/reconsiliation/internal/risk"
	risks "github.com/rafaeldepontes/reconsiliation/internal/risk/server"
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

	return mux
}
