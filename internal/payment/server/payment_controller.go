package server

import (
	"encoding/json"
	"net/http"

	"github.com/rafaeldepontes/goplo/internal/idempotency"
	"github.com/rafaeldepontes/goplo/internal/payment"
	ps "github.com/rafaeldepontes/goplo/internal/payment/service"
	"github.com/rafaeldepontes/goplo/internal/util"
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

	// Check cache looking for the key. (48h for TLS)
	paymentID, err := c.service.CheckKey(idempotencyKey)
	if err == nil {
		w.WriteHeader(204)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"payment_id": paymentID,
			"status":     "processed",
		})
		return
	}

	paymentID, err = c.service.ProcessPayment()
	if err != nil {
		util.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(204)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"payment_id": paymentID,
		"status":     "processed",
	})
}
