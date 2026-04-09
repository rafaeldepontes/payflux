package server

import (
	"encoding/json"
	"net/http"

	"github.com/rafaeldepontes/reconsiliation/internal/reconciliation"
	rs "github.com/rafaeldepontes/reconsiliation/internal/reconciliation/service"
	"github.com/rafaeldepontes/reconsiliation/internal/util"
)

type controller struct {
	service reconciliation.Service
}

func NewController() reconciliation.Controller {
	return &controller{
		service: rs.NewService(),
	}
}

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
