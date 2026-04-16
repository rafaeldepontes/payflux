package server

import (
	"encoding/json"
	"net/http"

	"github.com/rafaeldepontes/reconsiliation/internal/risk"
	"github.com/rafaeldepontes/reconsiliation/internal/util"
)

type controller struct {
	service risk.Service
}

func NewController(svc risk.Service) risk.Controller {
	return &controller{
		service: svc,
	}
}

// GetRiskEvaluation godoc
// @Summary Get risk evaluation
// @Description Returns the risk score and flags for a specific transaction
// @Tags risk
// @Produce  json
// @Param transaction_id path string true "Transaction ID"
// @Success 200 {object} model.RiskEvaluation
// @Failure 404 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Router /risk/{transaction_id} [get]
func (c *controller) GetRiskEvaluation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("transaction_id")
	res, err := c.service.GetResult(id)
	if err != nil {
		util.HandleError(w, "risk evaluation not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
