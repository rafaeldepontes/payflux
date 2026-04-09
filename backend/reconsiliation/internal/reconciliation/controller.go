package reconciliation

import "net/http"

type Controller interface {
	GetReconciliationResult(w http.ResponseWriter, r *http.Request)
	ListExceptions(w http.ResponseWriter, r *http.Request)
	CreateSettlementRecord(w http.ResponseWriter, r *http.Request)
}
