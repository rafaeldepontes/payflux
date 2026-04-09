package server

import (
	"encoding/json"
	"net/http"

	"github.com/rafaeldepontes/ledger/internal/idempotency"
	"github.com/rafaeldepontes/ledger/internal/payment"
	pm "github.com/rafaeldepontes/ledger/internal/payment/model"
	ps "github.com/rafaeldepontes/ledger/internal/payment/service"
	"github.com/rafaeldepontes/ledger/internal/util"
)

type controller struct {
	service payment.Service
}

func NewController() payment.Controller {
	return &controller{
		service: ps.NewService(),
	}
}

func (c controller) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		util.HandleError(w, "missing idempotency-key", http.StatusBadRequest)
		return
	}

	if len(idempotencyKey) != idempotency.IdempotencyKeySize {
		util.HandleError(w, "invalid idempotency-key", http.StatusBadRequest)
		return
	}

	// Check cache looking for the key. (48h for TTL)
	pres, err := c.service.CheckKey(idempotencyKey)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pres)
		return
	}

	var payment pm.PaymentReq
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		util.HandleError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := c.service.ProcessPayment(idempotencyKey, payment)
	if err != nil {
		util.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (c controller) GetPayment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		util.HandleError(w, "missing payment id", http.StatusBadRequest)
		return
	}

	res, err := c.service.GetPayment(id)
	if err != nil {
		util.HandleError(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (c controller) RefundPayment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		util.HandleError(w, "missing payment id", http.StatusBadRequest)
		return
	}

	var req pm.RefundReq
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.HandleError(w, "invalid request body", http.StatusBadRequest)
			return
		}
	}

	res, err := c.service.RefundPayment(id, req)
	if err != nil {
		util.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
