package server

import (
	"encoding/json"
	"net/http"

	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation"
	"github.com/rafaeldepontes/reconsiliation/internal/util"
)

type controller struct {
	service reconciliation.Service
}

func NewController(svc reconciliation.Service) reconciliation.Controller {
	return &controller{
		service: svc,
	}
}

// GetReconciliationResult godoc
// @Summary Get reconciliation result
// @Description Returns the reconciliation result for a transaction
// @Tags reconciliation
// @Produce json
// @Param transaction_id path string true "Transaction ID"
// @Success 200 {object} model.ReconciliationResult
// @Failure 404 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Router /reconciliation/{transaction_id} [get]
func (c *controller) GetReconciliationResult(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("transaction_id")
	res, err := c.service.GetResult(id)
	if err != nil {
		util.HandleError(w, "reconciliation result not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (c *controller) ListExceptions(w http.ResponseWriter, r *http.Request) {
	res, err := c.service.ListExceptions()
	if err != nil {
		util.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// CreateSettlementRecord godoc
// @Summary Create a settlement record
// @Description Creates or updates a settlement record for matching
// @Tags reconciliation
// @Accept  json
// @Produce  json
// @Param settlement body object true "Settlement Record"
// @Success 201
// @Failure 400 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Router /settlements [post]
func (c *controller) CreateSettlementRecord(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TransactionID string `json:"transaction_id"`
		Amount        int64  `json:"amount"`
		Status        string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.HandleError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.CreateSettlementRecord(req.TransactionID, req.Amount, req.Status); err != nil {
		util.HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
